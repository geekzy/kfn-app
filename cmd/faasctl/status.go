package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/geekzy/kfn-app/pkg/config"
	"github.com/geekzy/kfn-app/pkg/registry"
)

var statusCmd = &cobra.Command{
	Use:   "status [function-name]",
	Short: "Show function status",
	Long:  `Show the status of a deployed function.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	functionName := args[0]

	// Initialize configuration
	if err := config.Init(); err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	// Create registry client
	client, err := registry.NewClient(config.Get())
	if err != nil {
		return fmt.Errorf("failed to create registry client: %w", err)
	}
	defer client.Close()

	// Get function
	ctx := context.Background()
	fn, err := client.GetFunction(ctx, functionName)
	if err != nil {
		return fmt.Errorf("failed to get function: %w", err)
	}

	// Display function status
	fmt.Printf("Name: %s\n", fn.Name)
	fmt.Printf("Language: %s\n", fn.Language)
	fmt.Printf("Runtime Version: %s\n", fn.RuntimeVersion)
	fmt.Printf("Runtime Image: %s\n", fn.RuntimeImage)
	fmt.Printf("Memory: %d MB\n", fn.Memory)
	fmt.Printf("Timeout: %s\n", fn.Timeout)
	fmt.Printf("Handler: %s\n", fn.Handler)
	fmt.Printf("State: %s\n", fn.State)
	fmt.Printf("Created At: %s\n", fn.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated At: %s\n", fn.UpdatedAt.Format("2006-01-02 15:04:05"))

	return nil
}