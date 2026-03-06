package panel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/fatedier/frp/client/configmgmt"
	"github.com/fatedier/frp/client/http/model"
	"github.com/fatedier/frp/client/proxy"
	"github.com/fatedier/frp/pkg/util/jsonx"
)

// ServiceAccessor is the subset of client.Service that the panel needs.
type ServiceAccessor interface {
	GetConfigManager() configmgmt.ConfigManager
	GetAllProxyStatusForPanel() []*proxy.WorkingStatus
	GetConfigFilePath() string
	GetServerAddr() string
}

// Handler is the panel HTTP handler.
type Handler struct {
	auth       *AuthStore
	svc        ServiceAccessor
	serverAddr string
	startTime  time.Time
}

// NewHandler creates a new panel handler.
func NewHandler(svc ServiceAccessor) *Handler {
	cfgPath := svc.GetConfigFilePath()
	return &Handler{
		auth:      NewAuthStore(cfgPath),
		svc:       svc,
		startTime: time.Now(),
	}
}

// Register mounts all panel routes onto the given subrouter.
// The subrouter is already stripped of the /panel prefix.
func (h *Handler) Register(r *mux.Router) {
	// Public routes (no session needed)
	r.HandleFunc("/api/login", h.handleLogin).Methods(http.MethodPost)
	r.HandleFunc("/api/ping", h.handlePing).Methods(http.MethodGet)

	// All other /api/* require a valid session
	api := r.PathPrefix("/api").Subrouter()
	api.Use(h.sessionMiddleware)

	api.HandleFunc("/logout", h.handleLogout).Methods(http.MethodPost)
	api.HandleFunc("/password", h.handleChangePassword).Methods(http.MethodPost)
	api.HandleFunc("/info", h.handleInfo).Methods(http.MethodGet)

	// frpc status (inline, no proxy needed — direct access to service)
	api.HandleFunc("/status", h.handleStatus).Methods(http.MethodGet)

	// Config file
	api.HandleFunc("/config", h.handleGetConfig).Methods(http.MethodGet)
	api.HandleFunc("/config", h.handlePutConfig).Methods(http.MethodPut)

	// Proxy CRUD via store API
	api.HandleFunc("/proxies", h.handleListProxies).Methods(http.MethodGet)
	api.HandleFunc("/proxies", h.handleCreateProxy).Methods(http.MethodPost)
	api.HandleFunc("/proxies/{name}", h.handleGetProxy).Methods(http.MethodGet)
	api.HandleFunc("/proxies/{name}", h.handleUpdateProxy).Methods(http.MethodPut)
	api.HandleFunc("/proxies/{name}", h.handleDeleteProxy).Methods(http.MethodDelete)

	// Reload (wraps existing frpc reload logic)
	api.HandleFunc("/reload", h.handleReload).Methods(http.MethodPost)

	// Serve the SPA for all non-API paths
	r.PathPrefix("/").HandlerFunc(h.handleUI)
}

// ── Middleware ────────────────────────────────────────────────────────────────

func (h *Handler) sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := SessionFromRequest(r)
		if !ValidSession(token) {
			jsonErr(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ── Auth endpoints ────────────────────────────────────────────────────────────

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonErr(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if !h.auth.Verify(req.Username, req.Password) {
		jsonErr(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	token := CreateSession(req.Username)
	SetSessionCookie(w, token)
	jsonOK(w, map[string]string{"username": req.Username})
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	token := SessionFromRequest(r)
	DestroySession(token)
	ClearSessionCookie(w)
	jsonOK(w, map[string]string{"status": "ok"})
}

func (h *Handler) handlePing(w http.ResponseWriter, r *http.Request) {
	// Also usable to check if already authenticated
	token := SessionFromRequest(r)
	jsonOK(w, map[string]interface{}{
		"ok":            true,
		"authenticated": ValidSession(token),
	})
}

func (h *Handler) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Current string `json:"current"`
		New     string `json:"new"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonErr(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.auth.ChangePassword(req.Current, req.New); err != nil {
		jsonErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonOK(w, map[string]string{"status": "ok"})
}

// ── Info / Status ─────────────────────────────────────────────────────────────

func (h *Handler) handleInfo(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, map[string]interface{}{
		"uptime":     time.Since(h.startTime).Round(time.Second).String(),
		"configFile": h.svc.GetConfigFilePath(),
		"serverAddr": h.svc.GetServerAddr(),
	})
}

func (h *Handler) handleStatus(w http.ResponseWriter, r *http.Request) {
	statuses := h.svc.GetAllProxyStatusForPanel()
	res := make(model.StatusResp)
	for _, s := range statuses {
		psr := model.ProxyStatusResp{
			Name:   s.Name,
			Type:   s.Type,
			Status: s.Phase,
			Err:    s.Err,
		}
		base := s.Cfg.GetBaseConfig()
		if base.LocalPort != 0 {
			psr.LocalAddr = fmt.Sprintf("%s:%d", base.LocalIP, base.LocalPort)
		}
		psr.Plugin = base.Plugin.Type
		if s.Err == "" {
			psr.RemoteAddr = s.RemoteAddr
		}
		res[s.Type] = append(res[s.Type], psr)
	}
	jsonOK(w, res)
}

// ── Config ────────────────────────────────────────────────────────────────────

func (h *Handler) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	mgr := h.svc.GetConfigManager()
	content, err := mgr.ReadConfigFile()
	if err != nil {
		jsonErr(w, fmt.Sprintf("read config: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(content))
}

func (h *Handler) handlePutConfig(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		jsonErr(w, "read body error", http.StatusBadRequest)
		return
	}
	mgr := h.svc.GetConfigManager()
	if err := mgr.WriteConfigFile(body); err != nil {
		jsonErr(w, fmt.Sprintf("write config: %v", err), http.StatusInternalServerError)
		return
	}
	jsonOK(w, map[string]string{"status": "saved"})
}

// ── Reload ────────────────────────────────────────────────────────────────────

func (h *Handler) handleReload(w http.ResponseWriter, r *http.Request) {
	mgr := h.svc.GetConfigManager()
	if err := mgr.ReloadFromFile(false); err != nil {
		jsonErr(w, fmt.Sprintf("reload: %v", err), http.StatusInternalServerError)
		return
	}
	jsonOK(w, map[string]string{"status": "reloaded"})
}

// ── Proxy CRUD ────────────────────────────────────────────────────────────────

func (h *Handler) handleListProxies(w http.ResponseWriter, r *http.Request) {
	mgr := h.svc.GetConfigManager()
	proxies, err := mgr.ListStoreProxies()
	if err != nil {
		jsonErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := model.ProxyListResp{Proxies: make([]model.ProxyDefinition, 0, len(proxies))}
	for _, p := range proxies {
		def, err := model.ProxyDefinitionFromConfigurer(p)
		if err != nil {
			continue
		}
		resp.Proxies = append(resp.Proxies, def)
	}
	jsonOK(w, resp)
}

func (h *Handler) handleGetProxy(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	mgr := h.svc.GetConfigManager()
	p, err := mgr.GetStoreProxy(name)
	if err != nil {
		jsonErr(w, err.Error(), http.StatusNotFound)
		return
	}
	def, err := model.ProxyDefinitionFromConfigurer(p)
	if err != nil {
		jsonErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonOK(w, def)
}

func (h *Handler) handleCreateProxy(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		jsonErr(w, "read body error", http.StatusBadRequest)
		return
	}
	var def model.ProxyDefinition
	if err := jsonx.Unmarshal(body, &def); err != nil {
		jsonErr(w, fmt.Sprintf("parse JSON: %v", err), http.StatusBadRequest)
		return
	}
	if err := def.Validate("", false); err != nil {
		jsonErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	cfg, err := def.ToConfigurer()
	if err != nil {
		jsonErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	mgr := h.svc.GetConfigManager()
	created, err := mgr.CreateStoreProxy(cfg)
	if err != nil {
		jsonErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, _ := model.ProxyDefinitionFromConfigurer(created)
	w.WriteHeader(http.StatusCreated)
	jsonOK(w, resp)
}

func (h *Handler) handleUpdateProxy(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	body, err := io.ReadAll(r.Body)
	if err != nil {
		jsonErr(w, "read body error", http.StatusBadRequest)
		return
	}
	var def model.ProxyDefinition
	if err := jsonx.Unmarshal(body, &def); err != nil {
		jsonErr(w, fmt.Sprintf("parse JSON: %v", err), http.StatusBadRequest)
		return
	}
	if err := def.Validate(name, true); err != nil {
		jsonErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	cfg, err := def.ToConfigurer()
	if err != nil {
		jsonErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	mgr := h.svc.GetConfigManager()
	updated, err := mgr.UpdateStoreProxy(name, cfg)
	if err != nil {
		jsonErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, _ := model.ProxyDefinitionFromConfigurer(updated)
	jsonOK(w, resp)
}

func (h *Handler) handleDeleteProxy(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	mgr := h.svc.GetConfigManager()
	if err := mgr.DeleteStoreProxy(name); err != nil {
		jsonErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonOK(w, map[string]string{"status": "deleted"})
}

// ── UI ────────────────────────────────────────────────────────────────────────

func (h *Handler) handleUI(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		jsonErr(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Inline the entire SPA — no external files needed
	_, _ = fmt.Fprint(w, panelHTML)
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func jsonOK(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func jsonErr(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}


