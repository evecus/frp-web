package server

import (
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/server/proxy"
	"github.com/fatedier/frp/server/registry"
)

// GetClientRegistry exposes the client registry for the panel.
func (svr *Service) GetClientRegistry() *registry.ClientRegistry {
	return svr.clientRegistry
}

// KickClientByKey closes the control connection for a client identified by
// registry key (format: "user.clientID" or runID).
func (svr *Service) KickClientByKey(key string) {
	info, ok := svr.clientRegistry.GetByKey(key)
	if !ok || !info.Online || info.RunID == "" {
		return
	}
	if ctl, ok := svr.ctlManager.GetByID(info.RunID); ok {
		ctl.Close()
	}
}

// GetProxiesByRunID returns all proxies that belong to the client with the
// given runID, by iterating the proxy manager.
func (svr *Service) GetProxiesByRunID(runID string) []proxy.Proxy {
	if runID == "" {
		return nil
	}
	var result []proxy.Proxy
	svr.pxyManager.ForEach(func(_ string, pxy proxy.Proxy) {
		lm := pxy.GetLoginMsg()
		if lm != nil && lm.RunID == runID {
			result = append(result, pxy)
		}
	})
	return result
}

// ProxyGetPort returns the remote-listening port for TCP/UDP proxies (0 for
// all other types, e.g. HTTP/HTTPS which use domains instead).
func ProxyGetPort(pxy proxy.Proxy) int {
	cfg := pxy.GetConfigurer()
	if cfg == nil {
		return 0
	}
	switch c := cfg.(type) {
	case *v1.TCPProxyConfig:
		return c.RemotePort
	case *v1.UDPProxyConfig:
		return c.RemotePort
	}
	return 0
}
