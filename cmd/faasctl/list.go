package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/geekzy/kfn-app/pkg/config"
	"github.com/geekzy/kfn-app/pkg/registry"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List deployed functions",
	Long:  `List all deployed functions registered in the etcd registry.`,
	RunE:  runList,
}

func runList(cmd *cobra.Command, args []string) error {
	// Initialize configuration
	if err := config.Init(); err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	// Create registry client
	client, err := registry.NewClient(config.Get())
	if err != nil {
		return fmt.Errorf("failed to create registry client: %w", err)
	}
	defer func() { _ = client.Close() }()

	// List functions
	ctx := context.Background()
	functions, err := client.ListFunctions(ctx)
	if err != nil {
		return fmt.Errorf("failed to list functions: %w", err)
	}

	// Display functions
	if len(functions) == 0 {
		fmt.Println("No functions deployed.")
		return nil
	}

	fmt.Printf("NAME\tLANGUAGE\tVERSION\tSTATE\n")
	for _, fn := range functions {
		fmt.Printf("%s\t%s\t%s\t%s\n",
			fn.Name,
			fn.Language,
			fn.RuntimeVersion,
			fn.State)
	}

	return nil
}
