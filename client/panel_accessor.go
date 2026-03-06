package client

import (
	"github.com/fatedier/frp/client/configmgmt"
	"github.com/fatedier/frp/client/proxy"
)

// GetConfigManager returns a ConfigManager for the panel to use.
func (svr *Service) GetConfigManager() configmgmt.ConfigManager {
	return newServiceConfigManager(svr)
}

// GetAllProxyStatusForPanel returns all current proxy statuses.
func (svr *Service) GetAllProxyStatusForPanel() []*proxy.WorkingStatus {
	return svr.getAllProxyStatus()
}

// GetConfigFilePath returns the path of the config file in use.
func (svr *Service) GetConfigFilePath() string {
	return svr.configFilePath
}

// GetServerAddr returns the frps server address.
func (svr *Service) GetServerAddr() string {
	return svr.common.ServerAddr
}
