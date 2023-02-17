package main

import (
	"flag"
	"os"

	"github.com/lburgazzoli/camel-go/cmd/camel/run"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "camel",
		Short: "camel",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	rootCmd.AddCommand(run.NewRunCmd())

	fs := flag.NewFlagSet("", flag.PanicOnError)

	rootCmd.PersistentFlags().AddGoFlagSet(fs)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
