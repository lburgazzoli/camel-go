package run

import (
	"github.com/spf13/cobra"
)

func NewRunCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "run",
		Short: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return &cmd
}
