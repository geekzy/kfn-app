package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/geekzy/kfn-app/pkg/config"
	"github.com/geekzy/kfn-app/pkg/registry"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [function-name]",
	Short: "Delete a function",
	Long:  `Delete a deployed function from the registry.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDelete,
}

func runDelete(cmd *cobra.Command, args []string) error {
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
	// explicitly ignore error on close; deferred cleanup should not mask
	// the primary error path.  lint requires us to handle the return value.
	defer func() { _ = client.Close() }()

	// Delete function
	ctx := context.Background()
	if err := client.DeleteFunction(ctx, functionName); err != nil {
		return fmt.Errorf("failed to delete function: %w", err)
	}

	fmt.Printf("Function '%s' deleted successfully.\n", functionName)
	return nil
}
