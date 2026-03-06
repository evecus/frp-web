// Package panel provides a web management panel embedded directly into frpc.
// It mounts under /panel/ on frpc's existing webServer port.
package panel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// AuthStore manages panel login credentials, stored in a JSON file
// next to the frpc config file so they persist across restarts.
type AuthStore struct {
	mu       sync.RWMutex
	filePath string
	data     authData
}

type authData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// sessions maps token → expiry.
var (
	sessionMu sync.Mutex
	sessions  = map[string]time.Time{}
)

func NewAuthStore(configFilePath string) *AuthStore {
	dir := "."
	if configFilePath != "" {
		dir = filepath.Dir(configFilePath)
	}
	a := &AuthStore{
		filePath: filepath.Join(dir, "frpc-panel-auth.json"),
		data:     authData{Username: "admin", Password: "admin"},
	}
	a.load()
	return a
}

func (a *AuthStore) load() {
	data, err := os.ReadFile(a.filePath)
	if err != nil {
		// First run: save defaults
		_ = a.save()
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	_ = json.Unmarshal(data, &a.data)
}

func (a *AuthStore) save() error {
	a.mu.RLock()
	data, err := json.MarshalIndent(a.data, "", "  ")
	a.mu.RUnlock()
	if err != nil {
		return err
	}
	return os.WriteFile(a.filePath, data, 0600)
}

func (a *AuthStore) Verify(username, password string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return username == a.data.Username && password == a.data.Password
}

func (a *AuthStore) ChangePassword(current, newPass string) error {
	a.mu.Lock()
	if current != a.data.Password {
		a.mu.Unlock()
		return fmt.Errorf("current password incorrect")
	}
	if len(newPass) < 4 {
		a.mu.Unlock()
		return fmt.Errorf("password too short (minimum 4 characters)")
	}
	a.data.Password = newPass
	a.mu.Unlock()
	return a.save()
}

func (a *AuthStore) Username() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.data.Username
}

// CreateSession returns a new session token for a user.
func CreateSession(username string) string {
	token := fmt.Sprintf("%s-%d", username, time.Now().UnixNano())
	sessionMu.Lock()
	sessions[token] = time.Now().Add(24 * time.Hour)
	sessionMu.Unlock()
	return token
}

// ValidSession reports whether token is a valid, non-expired session.
func ValidSession(token string) bool {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	exp, ok := sessions[token]
	return ok && time.Now().Before(exp)
}

// DestroySession invalidates a session token.
func DestroySession(token string) {
	sessionMu.Lock()
	delete(sessions, token)
	sessionMu.Unlock()
}

// SessionFromRequest extracts the session token from a request cookie.
func SessionFromRequest(r *http.Request) string {
	c, err := r.Cookie("frpc_panel_session")
	if err != nil {
		return ""
	}
	return c.Value
}

// SetSessionCookie writes the session cookie to the response.
func SetSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "frpc_panel_session",
		Value:    token,
		Path:     "/panel/",
		HttpOnly: true,
		MaxAge:   86400,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearSessionCookie clears the session cookie.
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "frpc_panel_session",
		Path:   "/panel/",
		MaxAge: -1,
	})
}
