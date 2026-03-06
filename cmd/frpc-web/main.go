package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/samber/lo"

	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/config/source"       // ← 改这行
	v1 "github.com/fatedier/frp/pkg/config/v1"
	frplog "github.com/fatedier/frp/pkg/util/log"
	"github.com/fatedier/frp/pkg/policy/security"     // ← 改这行
)

// ─── Data types ───────────────────────────────────────────────────────────────

type ServerProfile struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ServerAddr string `json:"serverAddr"`
	ServerPort int    `json:"serverPort"`
	AuthToken  string `json:"authToken"`
	Username   string `json:"username"`
}

type Settings struct {
	PanelUser       string          `json:"panelUser"`
	PanelPass       string          `json:"panelPass"`
	ActiveProfileID string          `json:"activeProfileID"`
	Profiles        []ServerProfile `json:"profiles"`
	LogLevel        string          `json:"logLevel"`
}

func defaultSettings() Settings {
	return Settings{
		PanelUser:       "admin",
		PanelPass:       "admin",
		ActiveProfileID: "default",
		Profiles: []ServerProfile{{
			ID:         "default",
			Name:       "默认服务器",
			ServerPort: 7000,
		}},
		LogLevel: "info",
	}
}

type Proxy struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	LocalIP      string `json:"localIP"`
	LocalPort    int    `json:"localPort"`
	RemotePort   int    `json:"remotePort,omitempty"`
	CustomDomain string `json:"customDomain,omitempty"`
	Subdomain    string `json:"subdomain,omitempty"`
	Disabled     bool   `json:"disabled"`
}

// ─── App ──────────────────────────────────────────────────────────────────────

type App struct {
	dataDir string
	listen  string

	mu       sync.RWMutex
	settings Settings

	svcMu     sync.Mutex
	svc       *client.Service
	svcCancel context.CancelFunc
	running   bool
	startErr  string

	sessions sync.Map
}

func main() {
	listenFlag := flag.String("listen", "127.0.0.1:7777", "panel listen address")
	flag.Parse()

	exe, err := os.Executable()
	if err != nil {
		exe = "."
	}
	dataDir := filepath.Join(filepath.Dir(exe), "data")

	app := &App{dataDir: dataDir, listen: *listenFlag, settings: defaultSettings()}
	if err := app.init(); err != nil {
		fmt.Fprintf(os.Stderr, "init: %v\n", err)
		os.Exit(1)
	}

	app.mu.RLock()
	user := app.settings.PanelUser
	app.mu.RUnlock()
	fmt.Printf("frpc-web  http://%s   login: %s   data: %s\n", *listenFlag, user, dataDir)

	srv := &http.Server{Addr: *listenFlag, Handler: app.router()}
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		app.stopFrpc()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "server: %v\n", err)
		os.Exit(1)
	}
}

// ─── Init / persistence ───────────────────────────────────────────────────────

func (a *App) init() error {
	if err := os.MkdirAll(a.dataDir, 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", a.dataDir, err)
	}
	if b, err := os.ReadFile(a.settingsFile()); err == nil {
		a.mu.Lock()
		if err2 := json.Unmarshal(b, &a.settings); err2 == nil {
			// Migrate old single-server format
			if len(a.settings.Profiles) == 0 {
				var old struct {
					PanelUser  string `json:"panelUser"`
					PanelPass  string `json:"panelPass"`
					ServerAddr string `json:"serverAddr"`
					ServerPort int    `json:"serverPort"`
					AuthToken  string `json:"authToken"`
					LogLevel   string `json:"logLevel"`
				}
				_ = json.Unmarshal(b, &old)
				port := old.ServerPort
				if port == 0 {
					port = 7000
				}
				a.settings = Settings{
					PanelUser: old.PanelUser, PanelPass: old.PanelPass,
					ActiveProfileID: "default", LogLevel: old.LogLevel,
					Profiles: []ServerProfile{{
						ID: "default", Name: "默认服务器",
						ServerAddr: old.ServerAddr, ServerPort: port, AuthToken: old.AuthToken,
					}},
				}
			}
		}
		a.mu.Unlock()
	}
	return a.saveSettings()
}

func (a *App) settingsFile() string { return filepath.Join(a.dataDir, "settings.json") }
func (a *App) proxiesFile() string  { return filepath.Join(a.dataDir, "proxies.json") }

func (a *App) saveSettings() error {
	a.mu.RLock()
	b, err := json.MarshalIndent(a.settings, "", "  ")
	a.mu.RUnlock()
	if err != nil {
		return err
	}
	return os.WriteFile(a.settingsFile(), b, 0600)
}

func (a *App) activeProfile() (ServerProfile, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, p := range a.settings.Profiles {
		if p.ID == a.settings.ActiveProfileID {
			return p, true
		}
	}
	if len(a.settings.Profiles) > 0 {
		return a.settings.Profiles[0], true
	}
	return ServerProfile{}, false
}

func (a *App) loadProxies() ([]Proxy, error) {
	b, err := os.ReadFile(a.proxiesFile())
	if os.IsNotExist(err) {
		return []Proxy{}, nil
	}
	if err != nil {
		return nil, err
	}
	var out []Proxy
	if err := json.Unmarshal(b, &out); err != nil {
		return []Proxy{}, nil
	}
	return out, nil
}

func (a *App) saveProxies(proxies []Proxy) error {
	if proxies == nil {
		proxies = []Proxy{}
	}
	b, _ := json.MarshalIndent(proxies, "", "  ")
	return os.WriteFile(a.proxiesFile(), b, 0644)
}

// ─── frpc lifecycle ───────────────────────────────────────────────────────────

func (a *App) buildCommonCfg(p ServerProfile) *v1.ClientCommonConfig {
	a.mu.RLock()
	logLevel := a.settings.LogLevel
	a.mu.RUnlock()
	cfg := &v1.ClientCommonConfig{}
	cfg.ServerAddr = p.ServerAddr
	cfg.ServerPort = lo.If(p.ServerPort > 0, p.ServerPort).Else(7000)
	cfg.Auth.Method = "token"
	cfg.Auth.Token = p.AuthToken
	cfg.User = p.Username
	cfg.LoginFailExit = lo.ToPtr(false)
	cfg.Log.Level = logLevel
	cfg.Log.To = "console"
	_ = cfg.Complete()
	return cfg
}

func (a *App) buildProxyCfgs() ([]v1.ProxyConfigurer, error) {
	proxies, err := a.loadProxies()
	if err != nil {
		return nil, err
	}
	var cfgs []v1.ProxyConfigurer
	for _, p := range proxies {
		if p.Disabled {
			continue
		}
		if c, err := proxyToCfg(p); err == nil {
			cfgs = append(cfgs, c)
		}
	}
	return cfgs, nil
}

func proxyToCfg(p Proxy) (v1.ProxyConfigurer, error) {
	t := lo.If(p.Type != "", p.Type).Else("tcp")
	ip := lo.If(p.LocalIP != "", p.LocalIP).Else("127.0.0.1")
	base := v1.ProxyBaseConfig{Name: p.Name}
	switch t {
	case "tcp":
		c := &v1.TCPProxyConfig{ProxyBaseConfig: base}
		c.LocalIP = ip; c.LocalPort = p.LocalPort; c.RemotePort = p.RemotePort
		c.Complete(); return c, nil
	case "udp":
		c := &v1.UDPProxyConfig{ProxyBaseConfig: base}
		c.LocalIP = ip; c.LocalPort = p.LocalPort; c.RemotePort = p.RemotePort
		c.Complete(); return c, nil
	case "http":
		c := &v1.HTTPProxyConfig{ProxyBaseConfig: base}
		c.LocalIP = ip; c.LocalPort = p.LocalPort
		if p.CustomDomain != "" {
			c.CustomDomains = []string{p.CustomDomain}
		}
		if p.Subdomain != "" {
			c.SubDomain = p.Subdomain
		}
		c.Complete(); return c, nil
	case "https":
		c := &v1.HTTPSProxyConfig{ProxyBaseConfig: base}
		c.LocalIP = ip; c.LocalPort = p.LocalPort
		if p.CustomDomain != "" {
			c.CustomDomains = []string{p.CustomDomain}
		}
		if p.Subdomain != "" {
			c.SubDomain = p.Subdomain
		}
		c.Complete(); return c, nil
	default:
		return nil, fmt.Errorf("unknown type: %s", t)
	}
}

func (a *App) startFrpc() error {
	a.svcMu.Lock()
	defer a.svcMu.Unlock()
	if a.running {
		return fmt.Errorf("already running")
	}
	profile, ok := a.activeProfile()
	if !ok || strings.TrimSpace(profile.ServerAddr) == "" {
		return fmt.Errorf("请先配置服务器地址")
	}
	common := a.buildCommonCfg(profile)
	frplog.InitLogger(common.Log.To, common.Log.Level, int(common.Log.MaxDays), common.Log.DisablePrintColor)

	proxyCfgs, err := a.buildProxyCfgs()
	if err != nil {
		return fmt.Errorf("load proxies: %w", err)
	}
	cs := source.NewConfigSource()
	if err := cs.ReplaceAll(proxyCfgs, nil); err != nil {
		return fmt.Errorf("config source: %w", err)
	}
	svc, err := client.NewService(client.ServiceOptions{
		Common:                 common,
		ConfigSourceAggregator: source.NewAggregator(cs),
		UnsafeFeatures:         security.NewUnsafeFeatures(nil),
	})
	if err != nil {
		return fmt.Errorf("new service: %w", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.svcCancel = cancel
	a.svc = svc
	a.running = true
	a.startErr = ""
	go func() {
		if err := svc.Run(ctx); err != nil && ctx.Err() == nil {
			a.svcMu.Lock()
			a.startErr = err.Error()
			a.svcMu.Unlock()
		}
		a.svcMu.Lock()
		a.running = false
		a.svc = nil
		a.svcMu.Unlock()
	}()
	return nil
}

func (a *App) stopFrpc() {
	a.svcMu.Lock()
	defer a.svcMu.Unlock()
	if !a.running {
		return
	}
	a.svcCancel()
	for i := 0; i < 30; i++ {
		a.svcMu.Unlock()
		time.Sleep(100 * time.Millisecond)
		a.svcMu.Lock()
		if !a.running {
			return
		}
	}
}

func (a *App) applyProxies() {
	a.svcMu.Lock()
	svc, running := a.svc, a.running
	a.svcMu.Unlock()
	if !running || svc == nil {
		return
	}
	cfgs, _ := a.buildProxyCfgs()
	_ = svc.UpdateAllConfigurer(cfgs, nil)
}

// ─── Sessions ─────────────────────────────────────────────────────────────────

func (a *App) newSession() string {
	tok := fmt.Sprintf("t%d", time.Now().UnixNano())
	a.sessions.Store(tok, time.Now().Add(24*time.Hour))
	return tok
}

func (a *App) validSession(r *http.Request) bool {
	c, err := r.Cookie("frpcweb")
	if err != nil {
		return false
	}
	v, ok := a.sessions.Load(c.Value)
	return ok && time.Now().Before(v.(time.Time))
}

// ─── Router ───────────────────────────────────────────────────────────────────

func (a *App) router() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/api/login", a.hLogin).Methods("POST")
	r.HandleFunc("/api/ping", a.hPing).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !a.validSession(r) {
				jsonErr(w, "unauthorized", 401)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
	api.HandleFunc("/logout", a.hLogout).Methods("POST")
	api.HandleFunc("/status", a.hStatus).Methods("GET")
	api.HandleFunc("/settings", a.hGetSettings).Methods("GET")
	api.HandleFunc("/settings", a.hPutSettings).Methods("PUT")
	api.HandleFunc("/profiles", a.hListProfiles).Methods("GET")
	api.HandleFunc("/profiles", a.hCreateProfile).Methods("POST")
	api.HandleFunc("/profiles/{id}", a.hUpdateProfile).Methods("PUT")
	api.HandleFunc("/profiles/{id}", a.hDeleteProfile).Methods("DELETE")
	api.HandleFunc("/profiles/{id}/activate", a.hActivateProfile).Methods("POST")
	api.HandleFunc("/profiles/{id}/test", a.hTestProfile).Methods("POST")
	api.HandleFunc("/frpc/start", a.hStart).Methods("POST")
	api.HandleFunc("/frpc/stop", a.hStop).Methods("POST")
	api.HandleFunc("/frpc/restart", a.hRestart).Methods("POST")
	api.HandleFunc("/proxies", a.hListProxies).Methods("GET")
	api.HandleFunc("/proxies", a.hCreateProxy).Methods("POST")
	api.HandleFunc("/proxies/{name}", a.hUpdateProxy).Methods("PUT")
	api.HandleFunc("/proxies/{name}", a.hDeleteProxy).Methods("DELETE")
	api.HandleFunc("/password", a.hChangePassword).Methods("POST")

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			jsonErr(w, "not found", 404)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, uiHTML)
	})
	return r
}

// ─── Handlers ─────────────────────────────────────────────────────────────────

func (a *App) hPing(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, map[string]bool{"ok": true, "authenticated": a.validSession(r)})
}

func (a *App) hLogin(w http.ResponseWriter, r *http.Request) {
	var req struct{ Username, Password string }
	json.NewDecoder(r.Body).Decode(&req)
	a.mu.RLock()
	ok := req.Username == a.settings.PanelUser && req.Password == a.settings.PanelPass
	a.mu.RUnlock()
	if !ok {
		jsonErr(w, "用户名或密码错误", 401)
		return
	}
	tok := a.newSession()
	http.SetCookie(w, &http.Cookie{Name: "frpcweb", Value: tok, Path: "/", HttpOnly: true, MaxAge: 86400})
	jsonOK(w, map[string]string{"username": req.Username})
}

func (a *App) hLogout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("frpcweb"); err == nil {
		a.sessions.Delete(c.Value)
	}
	http.SetCookie(w, &http.Cookie{Name: "frpcweb", Path: "/", MaxAge: -1})
	jsonOK(w, nil)
}

func (a *App) hStatus(w http.ResponseWriter, r *http.Request) {
	a.svcMu.Lock()
	running, svc, errStr := a.running, a.svc, a.startErr
	a.svcMu.Unlock()
	profile, _ := a.activeProfile()
	res := map[string]any{
		"running": running, "error": errStr,
		"serverAddr": profile.ServerAddr, "serverPort": profile.ServerPort,
		"profileName": profile.Name, "username": profile.Username,
	}
	if running && svc != nil {
		all := svc.GetAllProxyStatusForPanel()
		list := make([]map[string]any, 0, len(all))
		for _, s := range all {
			base := s.Cfg.GetBaseConfig()
			p := map[string]any{"name": s.Name, "type": s.Type, "status": s.Phase, "err": s.Err}
			if base.LocalPort != 0 {
				p["localAddr"] = fmt.Sprintf("%s:%d", base.LocalIP, base.LocalPort)
			}
			if s.Err == "" {
				p["remoteAddr"] = s.RemoteAddr
			}
			list = append(list, p)
		}
		res["proxies"] = list
	}
	jsonOK(w, res)
}

func (a *App) hGetSettings(w http.ResponseWriter, r *http.Request) {
	a.mu.RLock()
	s := a.settings
	a.mu.RUnlock()
	s.PanelPass = ""
	jsonOK(w, s)
}

func (a *App) hPutSettings(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LogLevel string `json:"logLevel"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	a.mu.Lock()
	a.settings.LogLevel = req.LogLevel
	a.mu.Unlock()
	a.saveSettings()
	jsonOK(w, nil)
}

func (a *App) hListProfiles(w http.ResponseWriter, r *http.Request) {
	a.mu.RLock()
	profiles := make([]ServerProfile, len(a.settings.Profiles))
	copy(profiles, a.settings.Profiles)
	active := a.settings.ActiveProfileID
	a.mu.RUnlock()
	jsonOK(w, map[string]any{"profiles": profiles, "activeID": active})
}

func (a *App) hCreateProfile(w http.ResponseWriter, r *http.Request) {
	var p ServerProfile
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		jsonErr(w, "bad request", 400)
		return
	}
	if strings.TrimSpace(p.Name) == "" || strings.TrimSpace(p.ServerAddr) == "" {
		jsonErr(w, "name and serverAddr required", 400)
		return
	}
	p.ID = fmt.Sprintf("p%d", time.Now().UnixNano())
	if p.ServerPort == 0 {
		p.ServerPort = 7000
	}
	a.mu.Lock()
	a.settings.Profiles = append(a.settings.Profiles, p)
	a.mu.Unlock()
	a.saveSettings()
	jsonOK(w, p)
}

func (a *App) hUpdateProfile(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var p ServerProfile
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		jsonErr(w, "bad request", 400)
		return
	}
	p.ID = id
	a.mu.Lock()
	found := false
	for i, pr := range a.settings.Profiles {
		if pr.ID == id {
			if p.AuthToken == "" {
				p.AuthToken = pr.AuthToken
			}
			a.settings.Profiles[i] = p
			found = true
			break
		}
	}
	a.mu.Unlock()
	if !found {
		jsonErr(w, "not found", 404)
		return
	}
	a.saveSettings()
	jsonOK(w, nil)
}

func (a *App) hDeleteProfile(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	a.mu.Lock()
	var keep []ServerProfile
	for _, p := range a.settings.Profiles {
		if p.ID != id {
			keep = append(keep, p)
		}
	}
	a.settings.Profiles = keep
	if a.settings.ActiveProfileID == id && len(keep) > 0 {
		a.settings.ActiveProfileID = keep[0].ID
	}
	a.mu.Unlock()
	a.saveSettings()
	jsonOK(w, nil)
}

func (a *App) hActivateProfile(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	a.mu.Lock()
	found := false
	for _, p := range a.settings.Profiles {
		if p.ID == id {
			a.settings.ActiveProfileID = id
			found = true
			break
		}
	}
	a.mu.Unlock()
	if !found {
		jsonErr(w, "not found", 404)
		return
	}
	a.saveSettings()
	jsonOK(w, nil)
}

func (a *App) hTestProfile(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	a.mu.RLock()
	var profile ServerProfile
	found := false
	for _, p := range a.settings.Profiles {
		if p.ID == id {
			profile = p
			found = true
			break
		}
	}
	a.mu.RUnlock()
	if !found {
		jsonErr(w, "not found", 404)
		return
	}
	port := lo.If(profile.ServerPort > 0, profile.ServerPort).Else(7000)
	addr := fmt.Sprintf("%s:%d", profile.ServerAddr, port)
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		jsonOK(w, map[string]any{"ok": false, "msg": fmt.Sprintf("连接失败: %v", err)})
		return
	}
	conn.Close()
	jsonOK(w, map[string]any{"ok": true, "msg": fmt.Sprintf("连接成功 (%s)", addr)})
}

func (a *App) hStart(w http.ResponseWriter, r *http.Request) {
	if err := a.startFrpc(); err != nil {
		jsonErr(w, err.Error(), 400)
		return
	}
	jsonOK(w, nil)
}

func (a *App) hStop(w http.ResponseWriter, r *http.Request) {
	a.stopFrpc()
	jsonOK(w, nil)
}

func (a *App) hRestart(w http.ResponseWriter, r *http.Request) {
	go func() {
		a.stopFrpc()
		time.Sleep(200 * time.Millisecond)
		a.startFrpc()
	}()
	jsonOK(w, nil)
}

func (a *App) hListProxies(w http.ResponseWriter, r *http.Request) {
	proxies, err := a.loadProxies()
	if err != nil {
		jsonErr(w, err.Error(), 500)
		return
	}
	jsonOK(w, proxies)
}

func (a *App) hCreateProxy(w http.ResponseWriter, r *http.Request) {
	var p Proxy
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		jsonErr(w, "bad request", 400)
		return
	}
	if strings.TrimSpace(p.Name) == "" || p.LocalPort == 0 {
		jsonErr(w, "name and localPort required", 400)
		return
	}
	if p.LocalIP == "" {
		p.LocalIP = "127.0.0.1"
	}
	if p.Type == "" {
		p.Type = "tcp"
	}
	proxies, _ := a.loadProxies()
	for _, e := range proxies {
		if e.Name == p.Name {
			jsonErr(w, "name already exists", 409)
			return
		}
	}
	proxies = append(proxies, p)
	if err := a.saveProxies(proxies); err != nil {
		jsonErr(w, err.Error(), 500)
		return
	}
	a.applyProxies()
	jsonOK(w, p)
}

func (a *App) hUpdateProxy(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	var p Proxy
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		jsonErr(w, "bad request", 400)
		return
	}
	if p.Name == "" {
		p.Name = name
	}
	proxies, _ := a.loadProxies()
	found := false
	for i, e := range proxies {
		if e.Name == name {
			proxies[i] = p
			found = true
			break
		}
	}
	if !found {
		jsonErr(w, "not found", 404)
		return
	}
	if err := a.saveProxies(proxies); err != nil {
		jsonErr(w, err.Error(), 500)
		return
	}
	a.applyProxies()
	jsonOK(w, p)
}

func (a *App) hDeleteProxy(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	proxies, _ := a.loadProxies()
	var out []Proxy
	for _, p := range proxies {
		if p.Name != name {
			out = append(out, p)
		}
	}
	if err := a.saveProxies(out); err != nil {
		jsonErr(w, err.Error(), 500)
		return
	}
	a.applyProxies()
	jsonOK(w, nil)
}

func (a *App) hChangePassword(w http.ResponseWriter, r *http.Request) {
	var req struct{ Current, New string }
	json.NewDecoder(r.Body).Decode(&req)
	a.mu.Lock()
	if req.Current != a.settings.PanelPass {
		a.mu.Unlock()
		jsonErr(w, "当前密码错误", 403)
		return
	}
	if len(req.New) < 4 {
		a.mu.Unlock()
		jsonErr(w, "密码至少4位", 400)
		return
	}
	a.settings.PanelPass = req.New
	a.mu.Unlock()
	a.saveSettings()
	jsonOK(w, nil)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func jsonOK(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if v == nil {
		w.Write([]byte("{}"))
		return
	}
	json.NewEncoder(w).Encode(v)
}

func jsonErr(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
