package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var invokeCmd = &cobra.Command{
	Use:   "invoke [function-name]",
	Short: "Invoke a function",
	Long:  `Invoke a deployed function with optional payload.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runInvoke,
}

func runInvoke(cmd *cobra.Command, args []string) error {
	functionName := args[0]

	// TODO: Implement function invocation logic
	// This would typically involve making an HTTP request to the scheduler/API gateway

	fmt.Printf("Invoking function '%s'...\n", functionName)
	fmt.Println("TODO: Implement function invocation logic")

	return nil
}