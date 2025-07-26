package providers

import (
	"homelab-manager/internal/hosts/providers"
	"os"

	"gopkg.in/yaml.v3"
)

type YAMLProvider struct {
	Path string
}

type YAMLHostConfig struct {
	Hosts []providers.HostEntry `yaml:"host"`
}

func (p *YAMLProvider) GetHostEntries() ([]providers.HostEntry, error) {
	data, err := os.ReadFile(p.Path)
	if err != nil {
		return nil, err
	}

	var cfg YAMLHostConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return cfg.Hosts, nil
}
