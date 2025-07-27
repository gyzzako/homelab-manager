package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"homelab-manager/internal"
	"homelab-manager/internal/hosts"
	"homelab-manager/internal/hosts/git"
	"homelab-manager/internal/hosts/providers"
	configProvider "homelab-manager/internal/hosts/providers/config"
	urlProvider "homelab-manager/internal/hosts/providers/url"

	"homelab-manager/utils"
)

var (
	provider              string
	path                  string
	shouldPushHostsToRepo bool
	token                 string
)

var HostCmd = &cobra.Command{
	Use:     CMD_NAME_HOST,
	Short:   "Apply entries to the system's hosts file",
	Example: getCommandExample(),
	Run:     runCommand,
}

func init() {
	HostCmd.Flags().StringVarP(&provider, "provider", "p", "", "Data provider (required)")
	HostCmd.MarkFlagRequired("provider")

	HostCmd.Flags().StringVarP(&path, "path", "", "", "Path/URL to data (required)")
	HostCmd.MarkFlagRequired("path")

	HostCmd.Flags().StringVarP(&token, "token", "t", "", "Authentication token (optional)")
	HostCmd.MarkFlagRequired("token")

	HostCmd.Flags().BoolVar(&shouldPushHostsToRepo, "push", false, "Push host data to remote Git repository")
}

func runCommand(cmd *cobra.Command, args []string) {
	provider, hostProvider, err := getHostProvider(provider)
	if err != nil {
		fmt.Printf("❌ Failed to get provider: %v\n", err)
		os.Exit(1)
	}

	if err := validateParams(provider); err != nil {
		fmt.Printf("❌ Invalid param(s) : %v\n", err)
		os.Exit(1)
	}

	hostEntries := updateHosts(hostProvider)

	if shouldPushHostsToRepo {
		pushHostsDataToGit(hostEntries)
	}
}

func getCommandExample() string {
	template := `  {APP_NAME} {HOST} --config config.yaml
  {APP_NAME} {HOST} -c config.yaml`

	msg := utils.ReplaceMany(template, map[string]string{
		"{APP_NAME}": internal.APP_NAME,
		"{HOST}":     CMD_NAME_HOST,
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

	ymlConfig := &configProvider.YAMLProvider{Path: path}
	gitCfg, err := ymlConfig.GetGitConfig()
	if err != nil {
		fmt.Printf("❌ Failed to get git config: %v\n", err)
		os.Exit(1)
	}

	if err := git.PushToGit(hostEntries, gitCfg); err != nil {
		fmt.Printf("❌ Git push failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Pushed host data to git successfully.")
}

func getHostProvider(provider string) (providers.Provider, providers.HostProvider, error) {
	var defaultProvider providers.Provider
	switch strings.ToUpper(provider) {
	case string(providers.ProviderConfig):
		return providers.ProviderConfig, &configProvider.YAMLProvider{Path: path}, nil
	case string(providers.ProviderUrl):
		return providers.ProviderUrl, &urlProvider.URLProvider{URL: path, Token: token}, nil
	default:
		return defaultProvider, nil, errors.New("Provider " + provider + " not supported")
	}
}

func validateParams(provider providers.Provider) error {
	if provider == providers.ProviderUrl && shouldPushHostsToRepo {
		return fmt.Errorf("push param not supported with URL provider")
	}

	return nil
}
