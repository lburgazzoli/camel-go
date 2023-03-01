package run

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lburgazzoli/camel-go/pkg/core"
	"github.com/spf13/cobra"

	// import processors.
	_ "github.com/lburgazzoli/camel-go/pkg/components/log"
	_ "github.com/lburgazzoli/camel-go/pkg/components/timer"
	_ "github.com/lburgazzoli/camel-go/pkg/components/wasm"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/to"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/transform"
)

func NewRunCmd() *cobra.Command {
	type opts struct {
		Routes []string
	}

	var o opts

	cmd := cobra.Command{
		Use:   "run",
		Short: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			done := make(chan os.Signal, 1)
			signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

			camelContext := core.NewContext()

			for i := range o.Routes {
				file, err := os.Open(o.Routes[i])
				if err != nil {
					return err
				}

				if err := camelContext.LoadRoutes(cmd.Context(), file); err != nil {
					return err
				}
			}

			<-done

			return nil
		},
	}

	cmd.Flags().StringSliceVar(&o.Routes, "route", nil, "routes")
	_ = cmd.MarkFlagRequired("routes")

	return &cmd
}
