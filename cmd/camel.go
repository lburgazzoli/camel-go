package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "camel-go",
}

// Execute --
func Execute() {
	if err := rootCmd.Execute(); err != nil {
	}
}
