package run

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/lburgazzoli/camel-go/pkg/core"
	"github.com/spf13/cobra"

	// helper to include everything.
	_ "github.com/lburgazzoli/camel-go/pkg/components/kafka"
	_ "github.com/lburgazzoli/camel-go/pkg/components/log"
	_ "github.com/lburgazzoli/camel-go/pkg/components/mqtt"
	_ "github.com/lburgazzoli/camel-go/pkg/components/timer"
	_ "github.com/lburgazzoli/camel-go/pkg/components/wasm"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/process"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/route"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/set_body"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/set_header"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/to"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/transform"
)

func NewRunCmd() *cobra.Command {
	type opts struct {
		Routes      []string
		Configs     []string
		Development bool
	}

	var o opts

	cmd := cobra.Command{
		Use:   "run",
		Short: "run",
		RunE: func(cmd *cobra.Command, args []string) error {

			if o.Development {
				l, err := zap.NewDevelopment()
				if err != nil {
					return err
				}

				core.L = l
			} else {
				l, err := zap.NewProduction()
				if err != nil {
					return err
				}

				core.L = l
			}

			defer func() {
				if core.L != nil {
					_ = core.L.Sync()
				}
			}()

			done := make(chan os.Signal, 1)
			signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

			ctx := context.Background()
			camelContext := core.NewContext(core.L)

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
	cmd.Flags().BoolVar(&o.Development, "dev", false, "development")

	_ = cmd.MarkFlagRequired("routes")

	return &cmd
}
