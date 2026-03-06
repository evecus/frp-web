// frps-web: frps with built-in web management panel.
// Run ./frps-web — no config file needed.
// All data is stored in ./data/ relative to the binary location.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	modelmetrics "github.com/fatedier/frp/pkg/metrics"
	"github.com/fatedier/frp/pkg/metrics/mem"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	frplog "github.com/fatedier/frp/pkg/util/log"
	"github.com/fatedier/frp/server"
)

// ─── Settings (saved to data/settings.json) ───────────────────────────────────

type Settings struct {
	PanelUser      string `json:"panelUser"`
	PanelPass      string `json:"panelPass"`
	BindPort       int    `json:"bindPort"`
	AuthToken      string `json:"authToken"`
	LogLevel       string `json:"logLevel"`
	VhostHTTPPort  int    `json:"vhostHTTPPort"`
	VhostHTTPSPort int    `json:"vhostHTTPSPort"`
	SubDomainHost  string `json:"subDomainHost"`
}

func defaultSettings() Settings {
	return Settings{
		PanelUser: "admin",
		PanelPass: "admin",
		BindPort:  7000,
		LogLevel:  "info",
	}
}

// ─── In-memory log ring ────────────────────────────────────────────────────────

const maxLogEntries = 500

type LogEntry struct {
	Time  string `json:"time"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

type LogRing struct {
	mu      sync.Mutex
	entries []LogEntry
}

func (lb *LogRing) Add(level, msg string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.entries = append(lb.entries, LogEntry{
		Time:  time.Now().Format("2006-01-02 15:04:05"),
		Level: level,
		Msg:   msg,
	})
	if len(lb.entries) > maxLogEntries {
		lb.entries = lb.entries[len(lb.entries)-maxLogEntries:]
	}
}

func (lb *LogRing) Snapshot() []LogEntry {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	out := make([]LogEntry, len(lb.entries))
	copy(out, lb.entries)
	return out
}

// ─── Blocklist (saved to data/blocklist.json) ─────────────────────────────────

type Blocklist struct {
	mu      sync.RWMutex
	blocked map[string]struct{}
}

func newBlocklist() *Blocklist { return &Blocklist{blocked: make(map[string]struct{})} }

func (b *Blocklist) Block(key string) {
	b.mu.Lock(); b.blocked[key] = struct{}{}; b.mu.Unlock()
}
func (b *Blocklist) Unblock(key string) {
	b.mu.Lock(); delete(b.blocked, key); b.mu.Unlock()
}
func (b *Blocklist) Has(key string) bool {
	b.mu.RLock(); defer b.mu.RUnlock(); _, ok := b.blocked[key]; return ok
}
func (b *Blocklist) Keys() []string {
	b.mu.RLock(); defer b.mu.RUnlock()
	out := make([]string, 0, len(b.blocked))
	for k := range b.blocked { out = append(out, k) }
	return out
}

// ─── App ──────────────────────────────────────────────────────────────────────

type App struct {
	dataDir string
	listen  string

	mu       sync.RWMutex
	settings Settings

	svcMu     sync.Mutex
	svc       *server.Service
	svcCancel context.CancelFunc
	running   bool
	startErr  string
	startedAt time.Time

	logRing   LogRing
	blocklist *Blocklist
	sessions  sync.Map // token(string) → expiry(time.Time)
}

// ─── main ─────────────────────────────────────────────────────────────────────

func main() {
	listenAddr := flag.String("listen", "0.0.0.0:7500", "panel HTTP listen address")
	flag.Parse()

	exe, err := os.Executable()
	if err != nil {
		exe = "."
	}
	dataDir := filepath.Join(filepath.Dir(exe), "data")

	app := &App{
		dataDir:   dataDir,
		listen:    *listenAddr,
		settings:  defaultSettings(),
		blocklist: newBlocklist(),
	}
	if err := app.init(); err != nil {
		fmt.Fprintf(os.Stderr, "init error: %v\n", err)
		os.Exit(1)
	}

	app.mu.RLock()
	user, pass := app.settings.PanelUser, app.settings.PanelPass
	app.mu.RUnlock()

	fmt.Printf("┌──────────────────────────────────────────────┐\n")
	fmt.Printf("│  frps-web management panel                   │\n")
	fmt.Printf("├──────────────────────────────────────────────┤\n")
	fmt.Printf("│  URL:   http://%-29s│\n", *listenAddr)
	fmt.Printf("│  Login: %-36s│\n", user+" / "+pass)
	fmt.Printf("│  Data:  %-36s│\n", dataDir)
	fmt.Printf("└──────────────────────────────────────────────┘\n")

	srv := &http.Server{
		Addr:         *listenAddr,
		Handler:      app.buildRouter(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Graceful shutdown
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		fmt.Println("\nShutting down…")
		app.stopFrps()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx) //nolint:errcheck
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "panel server: %v\n", err)
		os.Exit(1)
	}
}

// ─── Init / persistence ────────────────────────────────────────────────────────

func (a *App) init() error {
	if err := os.MkdirAll(a.dataDir, 0o755); err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}
	// Load settings (ignore first-run missing file)
	if raw, err := os.ReadFile(a.settingsPath()); err == nil {
		a.mu.Lock()
		_ = json.Unmarshal(raw, &a.settings)
		a.mu.Unlock()
	}
	// Load blocklist
	if raw, err := os.ReadFile(a.blocklistPath()); err == nil {
		var keys []string
		if json.Unmarshal(raw, &keys) == nil {
			for _, k := range keys {
				a.blocklist.Block(k)
			}
		}
	}
	return a.saveSettings()
}

func (a *App) settingsPath() string  { return filepath.Join(a.dataDir, "settings.json") }
func (a *App) blocklistPath() string { return filepath.Join(a.dataDir, "blocklist.json") }

func (a *App) saveSettings() error {
	a.mu.RLock()
	raw, err := json.MarshalIndent(a.settings, "", "  ")
	a.mu.RUnlock()
	if err != nil {
		return err
	}
	return os.WriteFile(a.settingsPath(), raw, 0o600)
}

func (a *App) saveBlocklist() {
	raw, _ := json.MarshalIndent(a.blocklist.Keys(), "", "  ")
	_ = os.WriteFile(a.blocklistPath(), raw, 0o644)
}

// ─── frps lifecycle ────────────────────────────────────────────────────────────

func (a *App) buildFrpsCfg() *v1.ServerConfig {
	a.mu.RLock()
	s := a.settings
	a.mu.RUnlock()

	cfg := &v1.ServerConfig{}
	cfg.BindAddr = "0.0.0.0"
	cfg.BindPort = s.BindPort
	cfg.Auth.Method = "token"
	cfg.Auth.Token = s.AuthToken
	cfg.Log.Level = s.LogLevel
	cfg.Log.To = "console"
	cfg.VhostHTTPPort = s.VhostHTTPPort
	cfg.VhostHTTPSPort = s.VhostHTTPSPort
	cfg.SubDomainHost = s.SubDomainHost
	_ = cfg.Complete()
	return cfg
}

func (a *App) startFrps() error {
	a.svcMu.Lock()
	defer a.svcMu.Unlock()

	if a.running {
		return fmt.Errorf("already running")
	}

	cfg := a.buildFrpsCfg()
	frplog.InitLogger(cfg.Log.To, cfg.Log.Level, int(cfg.Log.MaxDays), cfg.Log.DisablePrintColor)

	// Enable in-memory metrics (normally gated behind webServer, we do it here).
	modelmetrics.EnableMem()

	svc, err := server.NewService(cfg)
	if err != nil {
		return fmt.Errorf("create frps service: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	a.svcCancel = cancel
	a.svc = svc
	a.running = true
	a.startErr = ""
	a.startedAt = time.Now()
	a.logRing.Add("info", fmt.Sprintf("frps started — listening on :%d", cfg.BindPort))

	go func() {
		svc.Run(ctx)
		a.svcMu.Lock()
		a.running = false
		a.svc = nil
		a.logRing.Add("info", "frps stopped")
		a.svcMu.Unlock()
	}()
	return nil
}

func (a *App) stopFrps() {
	a.svcMu.Lock()
	defer a.svcMu.Unlock()
	if !a.running {
		return
	}
	a.svcCancel()
	// Wait up to 3 s for goroutine to mark running=false
	for i := 0; i < 30; i++ {
		a.svcMu.Unlock()
		time.Sleep(100 * time.Millisecond)
		a.svcMu.Lock()
		if !a.running {
			return
		}
	}
}

// ─── Session management ────────────────────────────────────────────────────────

func (a *App) newSession() string {
	token := fmt.Sprintf("s%d", time.Now().UnixNano())
	a.sessions.Store(token, time.Now().Add(24*time.Hour))
	return token
}

func (a *App) validSession(r *http.Request) bool {
	c, err := r.Cookie("frpsweb")
	if err != nil {
		return false
	}
	v, ok := a.sessions.Load(c.Value)
	return ok && time.Now().Before(v.(time.Time))
}

// ─── Router ────────────────────────────────────────────────────────────────────

func (a *App) buildRouter() http.Handler {
	r := mux.NewRouter()

	// Public endpoints
	r.HandleFunc("/api/login", a.hLogin).Methods("POST")
	r.HandleFunc("/api/ping", a.hPing).Methods("GET")

	// Authenticated API group
	authd := r.PathPrefix("/api").Subrouter()
	authd.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !a.validSession(r) {
				apiErr(w, "unauthorized", 401)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	authd.HandleFunc("/logout", a.hLogout).Methods("POST")
	authd.HandleFunc("/password", a.hChangePassword).Methods("POST")

	authd.HandleFunc("/frps/status", a.hStatus).Methods("GET")
	authd.HandleFunc("/frps/start", a.hStart).Methods("POST")
	authd.HandleFunc("/frps/stop", a.hStop).Methods("POST")
	authd.HandleFunc("/frps/restart", a.hRestart).Methods("POST")

	authd.HandleFunc("/settings", a.hGetSettings).Methods("GET")
	authd.HandleFunc("/settings", a.hPutSettings).Methods("PUT")

	authd.HandleFunc("/clients", a.hClients).Methods("GET")
	authd.HandleFunc("/clients/{key}/block", a.hBlock).Methods("POST")
	authd.HandleFunc("/clients/{key}/unblock", a.hUnblock).Methods("POST")

	authd.HandleFunc("/proxies", a.hProxies).Methods("GET")
	authd.HandleFunc("/proxies/{name}/traffic", a.hProxyTraffic).Methods("GET")

	authd.HandleFunc("/logs", a.hLogs).Methods("GET")

	// SPA — everything else
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			apiErr(w, "not found", 404)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, uiHTML)
	})
	return r
}

// ─── Handlers ─────────────────────────────────────────────────────────────────

func (a *App) hPing(w http.ResponseWriter, r *http.Request) {
	apiOK(w, map[string]bool{"ok": true, "authenticated": a.validSession(r)})
}

func (a *App) hLogin(w http.ResponseWriter, r *http.Request) {
	var req struct{ Username, Password string }
	json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
	a.mu.RLock()
	ok := req.Username == a.settings.PanelUser && req.Password == a.settings.PanelPass
	a.mu.RUnlock()
	if !ok {
		a.logRing.Add("warn", fmt.Sprintf("failed login attempt for user %q", req.Username))
		apiErr(w, "invalid credentials", 401)
		return
	}
	token := a.newSession()
	http.SetCookie(w, &http.Cookie{
		Name: "frpsweb", Value: token, Path: "/",
		HttpOnly: true, SameSite: http.SameSiteLaxMode, MaxAge: 86400,
	})
	a.logRing.Add("info", fmt.Sprintf("panel login: %s", req.Username))
	apiOK(w, map[string]string{"username": req.Username})
}

func (a *App) hLogout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("frpsweb"); err == nil {
		a.sessions.Delete(c.Value)
	}
	http.SetCookie(w, &http.Cookie{Name: "frpsweb", Path: "/", MaxAge: -1})
	apiOK(w, nil)
}

func (a *App) hChangePassword(w http.ResponseWriter, r *http.Request) {
	var req struct{ Current, New string }
	json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
	a.mu.Lock()
	if req.Current != a.settings.PanelPass {
		a.mu.Unlock()
		apiErr(w, "wrong current password", 403)
		return
	}
	if len(strings.TrimSpace(req.New)) < 4 {
		a.mu.Unlock()
		apiErr(w, "new password too short (min 4 chars)", 400)
		return
	}
	a.settings.PanelPass = req.New
	a.mu.Unlock()
	a.saveSettings() //nolint:errcheck
	a.logRing.Add("info", "panel password changed")
	apiOK(w, nil)
}

func (a *App) hStatus(w http.ResponseWriter, r *http.Request) {
	a.svcMu.Lock()
	running, errStr, startedAt := a.running, a.startErr, a.startedAt
	a.svcMu.Unlock()
	a.mu.RLock()
	bindPort := a.settings.BindPort
	a.mu.RUnlock()

	res := map[string]any{
		"running":  running,
		"error":    errStr,
		"bindPort": bindPort,
	}
	if running && !startedAt.IsZero() {
		res["startTime"] = startedAt.Format("2006-01-02 15:04:05")
		res["uptime"] = time.Since(startedAt).Round(time.Second).String()
		stats := mem.StatsCollector.GetServer()
		res["totalTrafficIn"] = stats.TotalTrafficIn
		res["totalTrafficOut"] = stats.TotalTrafficOut
		res["curConns"] = stats.CurConns
		res["clientCounts"] = stats.ClientCounts
		res["proxyTypeCounts"] = stats.ProxyTypeCounts
	}
	apiOK(w, res)
}

func (a *App) hStart(w http.ResponseWriter, r *http.Request) {
	if err := a.startFrps(); err != nil {
		apiErr(w, err.Error(), 400)
		return
	}
	apiOK(w, nil)
}

func (a *App) hStop(w http.ResponseWriter, r *http.Request) {
	a.stopFrps()
	apiOK(w, nil)
}

func (a *App) hRestart(w http.ResponseWriter, r *http.Request) {
	go func() {
		a.stopFrps()
		time.Sleep(300 * time.Millisecond)
		if err := a.startFrps(); err != nil {
			a.svcMu.Lock()
			a.startErr = err.Error()
			a.svcMu.Unlock()
		}
	}()
	apiOK(w, nil)
}

func (a *App) hGetSettings(w http.ResponseWriter, r *http.Request) {
	a.mu.RLock()
	s := a.settings
	a.mu.RUnlock()
	s.PanelPass = "" // never send password over the wire
	apiOK(w, s)
}

func (a *App) hPutSettings(w http.ResponseWriter, r *http.Request) {
	var req Settings
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiErr(w, "bad request body", 400)
		return
	}
	if req.BindPort <= 0 || req.BindPort > 65535 {
		apiErr(w, "bindPort must be 1-65535", 400)
		return
	}
	a.mu.Lock()
	a.settings.BindPort = req.BindPort
	a.settings.AuthToken = req.AuthToken
	a.settings.LogLevel = req.LogLevel
	a.settings.VhostHTTPPort = req.VhostHTTPPort
	a.settings.VhostHTTPSPort = req.VhostHTTPSPort
	a.settings.SubDomainHost = req.SubDomainHost
	a.mu.Unlock()
	if err := a.saveSettings(); err != nil {
		apiErr(w, err.Error(), 500)
		return
	}
	a.logRing.Add("info", fmt.Sprintf("settings saved — bindPort=%d", req.BindPort))
	apiOK(w, nil)
}

func (a *App) hClients(w http.ResponseWriter, r *http.Request) {
	a.svcMu.Lock()
	svc := a.svc
	running := a.running
	a.svcMu.Unlock()

	if !running || svc == nil {
		apiOK(w, []any{})
		return
	}

	clients := svc.GetClientRegistry().List()
	out := make([]map[string]any, 0, len(clients))
	for _, c := range clients {
		// Collect occupied remote ports
		proxies := svc.GetProxiesByRunID(c.RunID)
		ports := make([]int, 0, len(proxies))
		for _, p := range proxies {
			if port := server.ProxyGetPort(p); port > 0 {
				ports = append(ports, port)
			}
		}

		out = append(out, map[string]any{
			"key":            c.Key,
			"user":           c.User,
			"clientID":       c.ClientID(),
			"runID":          c.RunID,
			"hostname":       c.Hostname,
			"ip":             c.IP,
			"version":        c.Version,
			"online":         c.Online,
			"firstConnected": fmtTime(c.FirstConnectedAt),
			"lastConnected":  fmtTime(c.LastConnectedAt),
			"disconnectedAt": fmtOptTime(c.DisconnectedAt),
			"blocked":        a.blocklist.Has(c.Key),
			"ports":          ports,
		})
	}
	apiOK(w, out)
}

func (a *App) hBlock(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	if key == "" {
		apiErr(w, "missing key", 400)
		return
	}
	a.blocklist.Block(key)
	a.saveBlocklist()
	a.logRing.Add("warn", fmt.Sprintf("blocked client: %s", key))
	// Immediately kick if online
	a.svcMu.Lock()
	svc := a.svc
	a.svcMu.Unlock()
	if svc != nil {
		svc.KickClientByKey(key)
	}
	apiOK(w, nil)
}

func (a *App) hUnblock(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	if key == "" {
		apiErr(w, "missing key", 400)
		return
	}
	a.blocklist.Unblock(key)
	a.saveBlocklist()
	a.logRing.Add("info", fmt.Sprintf("unblocked client: %s", key))
	apiOK(w, nil)
}

func (a *App) hProxies(w http.ResponseWriter, r *http.Request) {
	a.svcMu.Lock()
	running := a.running
	a.svcMu.Unlock()
	if !running {
		apiOK(w, []any{})
		return
	}

	types := []string{"tcp", "udp", "http", "https", "stcp", "xtcp", "tcpmux", "sudp"}
	all := make([]map[string]any, 0)
	for _, t := range types {
		for _, ps := range mem.StatsCollector.GetProxiesByType(t) {
			all = append(all, map[string]any{
				"name":            ps.Name,
				"type":            ps.Type,
				"user":            ps.User,
				"clientID":        ps.ClientID,
				"todayTrafficIn":  ps.TodayTrafficIn,
				"todayTrafficOut": ps.TodayTrafficOut,
				"curConns":        ps.CurConns,
				"lastStartTime":   ps.LastStartTime,
				"lastCloseTime":   ps.LastCloseTime,
			})
		}
	}
	apiOK(w, all)
}

func (a *App) hProxyTraffic(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	ti := mem.StatsCollector.GetProxyTraffic(name)
	if ti == nil {
		apiOK(w, map[string]any{"name": name, "trafficIn": []int64{}, "trafficOut": []int64{}})
		return
	}
	apiOK(w, map[string]any{
		"name":       ti.Name,
		"trafficIn":  ti.TrafficIn,
		"trafficOut": ti.TrafficOut,
	})
}

func (a *App) hLogs(w http.ResponseWriter, r *http.Request) {
	apiOK(w, a.logRing.Snapshot())
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func apiOK(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if v == nil {
		w.Write([]byte("{}")) //nolint:errcheck
		return
	}
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}

func apiErr(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg}) //nolint:errcheck
}

func fmtTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func fmtOptTime(t time.Time) string { return fmtTime(t) }
