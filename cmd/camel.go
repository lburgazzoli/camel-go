package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "desertship",
}

// Execute --
func Execute() {
	if err := rootCmd.Execute(); err != nil {
	}
}
