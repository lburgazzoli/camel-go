package run

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lburgazzoli/camel-go/pkg/core"
	"github.com/spf13/cobra"
)

func NewRunCmd() *cobra.Command {
	type opts struct {
		Routes  []string
		Configs []string
	}

	var o opts

	cmd := cobra.Command{
		Use:   "run",
		Short: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			done := make(chan os.Signal, 1)
			signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

			ctx := context.Background()
			camelContext := core.NewContext()

			for i := range o.Configs {
				if err := camelContext.Properties().AddSource(o.Configs[i]); err != nil {
					return err
				}
			}

			for i := range o.Routes {
				file, err := os.Open(o.Routes[i])
				if err != nil {
					return err
				}

				if err := camelContext.LoadRoutes(cmd.Context(), file); err != nil {
					return err
				}
			}

			if err := camelContext.Start(ctx); err != nil {
				return err
			}

			<-done

			return camelContext.Stop(ctx)
		},
	}

	cmd.Flags().StringSliceVar(&o.Routes, "route", nil, "routes")
	cmd.Flags().StringSliceVar(&o.Configs, "config", nil, "configs")

	_ = cmd.MarkFlagRequired("routes")

	return &cmd
}
