package hosts

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"homelab-manager/internal"
	"homelab-manager/internal/hosts"
	hostConfig "homelab-manager/internal/hosts/config"
	"homelab-manager/internal/hosts/git"
	"homelab-manager/internal/hosts/providers"
	configProvider "homelab-manager/internal/hosts/providers/config"
	sqlProvider "homelab-manager/internal/hosts/providers/sql"
	urlProvider "homelab-manager/internal/hosts/providers/url"

	"homelab-manager/utils"
)

const (
	CMD_NAME_HOST       = "host"
	PARAM_NAME_CONFIG   = "config"
	PARAM_NAME_PROVIDER = "provider"
	PARAM_NAME_PATH     = "path"
	PARAM_NAME_TOKEN    = "token"
	PARAM_NAME_PUSH     = "push"
	PARAM_NAME_TYPE     = "type"
	PARAM_NAME_QUERY    = "query"
)

var (
	config                hostConfig.HostConfig
	configPath            string
	provider              string
	path                  string
	shouldPushHostsToRepo bool
	token                 string
	dbType                string
	query                 string
)

var HostCmd = &cobra.Command{
	Use:     CMD_NAME_HOST,
	Short:   "Apply entries to the system's hosts file",
	Example: getCommandExample(),
	Run:     runCommand,
}

func init() {
	HostCmd.Flags().StringVarP(&configPath, PARAM_NAME_CONFIG, "c", "", "Config file path")
	HostCmd.Flags().StringVarP(&provider, PARAM_NAME_PROVIDER, "p", "", "Data provider")
	HostCmd.Flags().StringVarP(&path, PARAM_NAME_PATH, "", "", "Path/URL to data")
	HostCmd.Flags().StringVarP(&token, PARAM_NAME_TOKEN, "", "", "Authentication token")
	HostCmd.Flags().BoolVar(&shouldPushHostsToRepo, PARAM_NAME_PUSH, false, "Push host data to remote Git repository")
	HostCmd.Flags().StringVarP(&dbType, PARAM_NAME_TYPE, "", "", "Database type")
	HostCmd.Flags().StringVarP(&query, PARAM_NAME_QUERY, "", "", "SQL query")
}

func runCommand(cmd *cobra.Command, args []string) {
	if err := loadHostConfig(); err != nil {
		fmt.Printf("❌ Failed to load config: %v\n", err)
		os.Exit(1)
	}

	providerType, hostProvider, err := getHostProvider()
	if err != nil {
		fmt.Printf("❌ Failed to get provider: %v\n", err)
		os.Exit(1)
	}

	if err := validateParams(providerType); err != nil {
		fmt.Printf("❌ Invalid parameter(s): %v\n", err)
		os.Exit(1)
	}

	hostEntries := updateHosts(hostProvider)

	if getOverriddenParam(&config.GitConfig.ShouldPush, &shouldPushHostsToRepo) {
		pushHostsDataToGit(hostEntries)
	}
}

func loadHostConfig() error {
	if len(configPath) == 0 {
		return nil
	}

	if err := hostConfig.LoadConfig(configPath, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	return nil
}

func getCommandExample() string {
	template := `  # Using CLI parameters
  {APP_NAME} {HOST} --{PARAM_NAME_PROVIDER} config --{PARAM_NAME_CONFIG} ./config.yml
  {APP_NAME} {HOST} --{PARAM_NAME_PROVIDER} url --{PARAM_NAME_PATH} https://api.example.com/hosts --{PARAM_NAME_TOKEN} your-auth-token
  {APP_NAME} {HOST} --{PARAM_NAME_PROVIDER} sql --{PARAM_NAME_PATH} ./sqlite.db  --{PARAM_NAME_TYPE} sqlite --{PARAM_NAME_QUERY} "select * from hosts"

  # Using config file
  {APP_NAME} {HOST} --{PARAM_NAME_CONFIG} ./config.yml
  
  # Using config file with CLI override
  {APP_NAME} {HOST} --{PARAM_NAME_CONFIG} ./config.yml --{PARAM_NAME_PROVIDER} url`

	msg := utils.ReplaceMany(template, map[string]string{
		"{APP_NAME}":            internal.APP_NAME,
		"{HOST}":                CMD_NAME_HOST,
		"{PARAM_NAME_PROVIDER}": PARAM_NAME_PROVIDER,
		"{PARAM_NAME_PATH}":     PARAM_NAME_PATH,
		"{PARAM_NAME_TOKEN}":    PARAM_NAME_TOKEN,
		"{PARAM_NAME_CONFIG}":   PARAM_NAME_CONFIG,
		"{PARAM_NAME_TYPE}":     PARAM_NAME_TYPE,
		"{PARAM_NAME_QUERY}":    PARAM_NAME_QUERY,
	})

	return msg
}

func updateHosts(provider providers.HostProvider) []providers.HostEntry {
	hostEntries, err := hosts.UpdateHosts(provider)
	if err != nil {
		fmt.Printf("❌ Failed to update hosts file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Hosts file updated successfully.")
	return hostEntries
}

func pushHostsDataToGit(hostEntries []providers.HostEntry) {
	fmt.Println("📦 Pushing host data to Git repository...")

	if err := git.PushToGit(hostEntries, config.GitConfig); err != nil {
		fmt.Printf("❌ Git push failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Pushed host data to git successfully.")
}

func getHostProvider() (providers.Provider, providers.HostProvider, error) {
	stringConfigProvider := string(config.HostProviderConfig.Provider)
	stringProvider := getOverriddenParam(&stringConfigProvider, &provider)
	var defaultProvider providers.Provider

	switch strings.ToUpper(stringProvider) {
	case string(providers.ProviderConfig):
		return providers.ProviderConfig, &configProvider.HostConfigProvider{HostConfig: config}, nil
	case string(providers.ProviderUrl):
		return providers.ProviderUrl,
			&urlProvider.URLProvider{
				URL:   getOverriddenParam(&config.HostProviderConfig.UrlParams.URL, &path),
				Token: getOverriddenParam(&config.HostProviderConfig.UrlParams.Token, &token),
			},
			nil
	case string(providers.ProviderSql):
		return providers.ProviderSql,
			&sqlProvider.SQLProvider{
				DataSource: getOverriddenParam(&config.HostProviderConfig.SqlParams.DataSource, &path),
				Type:       getOverriddenParam(&config.HostProviderConfig.SqlParams.Type, &path),
				Query:      getOverriddenParam(&config.HostProviderConfig.SqlParams.Query, &path),
			}, nil
	default:
		return defaultProvider, nil, errors.New("Provider " + provider + " not supported")
	}
}

func validateParams(providerType providers.Provider) error {
	if providerType == providers.ProviderConfig && len(config.HostEntries) == 0 {
		return fmt.Errorf("hosts data is missing")
	}

	if providerType == providers.ProviderUrl {
		if utils.IsEmpty(getOverriddenParam(&config.HostProviderConfig.UrlParams.URL, &path)) {
			return fmt.Errorf("URL is missing")
		}
	}

	if providerType == providers.ProviderSql {
		if utils.IsEmpty(getOverriddenParam(&config.HostProviderConfig.SqlParams.DataSource, &path)) {
			return fmt.Errorf("data source is missing")
		}
		if utils.IsEmpty(getOverriddenParam(&config.HostProviderConfig.SqlParams.Type, &dbType)) {
			return fmt.Errorf("database type is missing")
		}
		if utils.IsEmpty(getOverriddenParam(&config.HostProviderConfig.SqlParams.Query, &query)) {
			return fmt.Errorf("SQL query is missing")
		}
	}

	if getOverriddenParam(&config.GitConfig.ShouldPush, &shouldPushHostsToRepo) {
		if utils.IsEmpty(config.GitConfig.URL) {
			return fmt.Errorf("git URL is missing")
		}
	}

	return nil
}

func getOverriddenParam[T any](configParam *T, cliParam *T) T {
	if configParam != nil {
		return *configParam
	}

	return *cliParam
}
