package providers

import (
	hostConfig "homelab-manager/internal/hosts/config"
	"homelab-manager/internal/hosts/providers"
)

type HostConfigProvider struct {
	HostConfig hostConfig.HostConfig
}

func (p *HostConfigProvider) GetHostEntries() ([]providers.HostEntry, error) {
	return p.HostConfig.HostEntries, nil
}
