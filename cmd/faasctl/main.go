package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version information
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "faasctl",
		Short: "KFn FaaS Platform CLI",
		Long: `faasctl is the command line interface for the KFn FaaS platform.
It allows you to deploy, manage, and invoke serverless functions.`,
	}

	// Add version command
	rootCmd.AddCommand(versionCmd)

	// Add subcommands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(invokeCmd)
	rootCmd.AddCommand(deleteCmd)

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
