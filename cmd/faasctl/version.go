package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of faasctl",
	Long:  `All software has versions. This is faasctl's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("faasctl version %s (built at %s)\n", Version, BuildTime)
	},
}
