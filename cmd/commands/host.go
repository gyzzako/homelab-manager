package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"homelab-manager/internal"
	"homelab-manager/internal/hosts"
	"homelab-manager/internal/hosts/providers"
	configProvider "homelab-manager/internal/hosts/providers/config"

	"homelab-manager/utils"
)

var (
	provider string
	path     string
)

var HostCmd = &cobra.Command{
	Use:     CMD_NAME_HOST,
	Short:   "Apply entries to the system's hosts file",
	Example: getCommandExample(),
	Run: func(cmd *cobra.Command, args []string) {

		provider, err := getHostProvider(provider, path)
		if err != nil {
			fmt.Printf("❌ Failed to get provider: %v\n", err)
			os.Exit(1)
		}

		if err := hosts.UpdateHosts(provider); err != nil {
			fmt.Printf("❌ Failed to update hosts file: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Hosts file updated successfully.")
	},
}

func init() {
	HostCmd.Flags().StringVarP(&provider, "provider", "p", "", "Data provider (required)")
	HostCmd.MarkFlagRequired("provider")

	HostCmd.Flags().StringVarP(&path, "path", "", "", "Path/URL to data (required)")
	HostCmd.MarkFlagRequired("path")
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

func getHostProvider(provider string, path string) (providers.HostProvider, error) {
	switch provider {
	case "config":
		return &configProvider.YAMLProvider{Path: path}, nil
	default:
		return nil, errors.New("Provider " + provider + " not supported")
	}
}
