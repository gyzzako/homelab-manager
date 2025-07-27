package main

import (
	"fmt"
	"homelab-manager/cmd/commands/hosts"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "homepage-manager",
		Short: "CLI tool for homelab management",
	}

	rootCmd.AddCommand(hosts.HostCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
