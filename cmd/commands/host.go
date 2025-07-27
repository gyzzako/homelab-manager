package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"homelab-manager/internal"
	"homelab-manager/internal/hosts"
	"homelab-manager/internal/hosts/git"
	"homelab-manager/internal/hosts/providers"
	configProvider "homelab-manager/internal/hosts/providers/config"

	"homelab-manager/utils"
)

var (
	provider              string
	path                  string
	shouldPushHostsToRepo bool
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

	HostCmd.Flags().BoolVar(&shouldPushHostsToRepo, "push", false, "Push host data to remote Git repository")
}

func runCommand(cmd *cobra.Command, args []string) {
	provider, err := getHostProvider(provider, path)
	if err != nil {
		fmt.Printf("❌ Failed to get provider: %v\n", err)
		os.Exit(1)
	}

	hostEntries := updateHosts(provider)

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

func getHostProvider(provider string, path string) (providers.HostProvider, error) {
	switch provider {
	case "config":
		return &configProvider.YAMLProvider{Path: path}, nil
	default:
		return nil, errors.New("Provider " + provider + " not supported")
	}
}
