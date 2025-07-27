package config

import (
	"fmt"
	"homelab-manager/internal/hosts/providers"
	"homelab-manager/internal/hosts/providers/sql"
	"homelab-manager/internal/hosts/providers/url"
	"os"

	"gopkg.in/yaml.v3"
)

type HostConfig struct {
	HostProviderConfig HostProviderConfig    `yaml:"provider"`
	GitConfig          GitConfig             `yaml:"git"`
	HostEntries        []providers.HostEntry `yaml:"data"`
}

type HostProviderConfig struct {
	Provider  providers.Provider `yaml:"type"`
	UrlParams url.URLProvider    `yaml:"url-params"`
	SqlParams sql.SQLProvider    `yaml:"sql-params"`
}

type GitConfig struct {
	ShouldPush bool   `yaml:"push"`
	URL        string `yaml:"url"`
	Token      string `yaml:"token"`
}

func LoadConfig(path string, cfg *HostConfig) error {
	configData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	if err := yaml.Unmarshal(configData, &cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	return nil
}
