package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/geekzy/kfn-app/pkg/config"
	"github.com/geekzy/kfn-app/pkg/registry"
)

var deployCmd = &cobra.Command{
	Use:   "deploy [function-name]",
	Short: "Deploy a function",
	Long:  `Deploy a function to the KFn platform.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDeploy,
}

func runDeploy(cmd *cobra.Command, args []string) error {
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
	// Close returns an error; log but don’t override any existing error.
	defer func() { _ = client.Close() }()

	// TODO: Parse function metadata from command line flags or config file
	// For now, we'll create a minimal function metadata
	fn := &registry.FunctionMetadata{
		Name:           functionName,
		Language:       "nodejs",
		RuntimeVersion: "18",
		RuntimeImage:   "node:18-alpine",
		Memory:         128,
		Timeout:        30 * time.Second,
		Handler:        "index.handler",
		State:          registry.FunctionPending,
	}

	// Deploy function
	ctx := context.Background()
	if err := client.CreateFunction(ctx, fn); err != nil {
		return fmt.Errorf("failed to deploy function: %w", err)
	}

	fmt.Printf("Function '%s' deployed successfully.\n", functionName)
	return nil
}
