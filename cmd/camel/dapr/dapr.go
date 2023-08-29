package dapr

import (
	"github.com/lburgazzoli/camel-go/cmd/camel/dapr/pub"
	"github.com/spf13/cobra"
)

func NewDaprCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "dapr",
		Short: "dapr",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.AddCommand(pub.NewPubCmd())

	return &cmd
}
