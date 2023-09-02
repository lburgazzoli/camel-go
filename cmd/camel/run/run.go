package run

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/lburgazzoli/camel-go/pkg/health"

	"github.com/lburgazzoli/camel-go/pkg/core"
	"github.com/spf13/cobra"

	// helper to include everything.
	_ "github.com/lburgazzoli/camel-go/pkg/components/dapr/pubsub"
	_ "github.com/lburgazzoli/camel-go/pkg/components/http"
	_ "github.com/lburgazzoli/camel-go/pkg/components/kafka"
	_ "github.com/lburgazzoli/camel-go/pkg/components/log"
	_ "github.com/lburgazzoli/camel-go/pkg/components/mqtt/v3"
	_ "github.com/lburgazzoli/camel-go/pkg/components/mqtt/v5"
	_ "github.com/lburgazzoli/camel-go/pkg/components/timer"
	_ "github.com/lburgazzoli/camel-go/pkg/components/wasm"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/choice"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/process"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/route"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/set_body"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/set_header"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/to"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/transform"
)

func NewRunCmd() *cobra.Command {
	type opts struct {
		Routes        []string
		Configs       []string
		Development   bool
		Health        bool
		HealthAddress string
		HealthPrefix  string
	}

	var o opts
	o.Health = true
	o.HealthPrefix = health.DefaultPrefix
	o.HealthAddress = health.DefaultAddress

	cmd := cobra.Command{
		Use:   "run",
		Short: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			var logger *slog.Logger

			if o.Development {
				logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
			} else {
				logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
			}

			slog.SetDefault(logger)

			done := make(chan os.Signal, 1)
			signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

			var h *health.Service

			if o.Health {
				h = health.New(o.HealthAddress, o.HealthPrefix, logger)

				if err := h.Start(cmd.Context()); err != nil {
					return err
				}
			}

			camelContext := core.NewContext(logger)

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
					_ = file.Close()

					return err
				}

				if err := file.Close(); err != nil {
					return err
				}
			}

			if err := camelContext.Start(cmd.Context()); err != nil {
				return err
			}

			<-done

			if err := camelContext.Stop(cmd.Context()); err != nil {
				return err
			}

			if h != nil {
				if err := h.Stop(cmd.Context()); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringSliceVar(&o.Routes, "route", nil, "routes")
	cmd.Flags().StringSliceVar(&o.Configs, "config", nil, "configs")
	cmd.Flags().BoolVar(&o.Development, "dev", false, "development")

	cmd.Flags().BoolVar(&o.Health, "health-check-enabled", o.Health, "health-check-enabled")
	cmd.Flags().StringVar(&o.HealthPrefix, "health-check-prefix", o.HealthPrefix, "health-check-prefix")
	cmd.Flags().StringVar(&o.HealthAddress, "health-check-address", o.HealthAddress, "health-check-address")

	_ = cmd.MarkFlagRequired("routes")

	return &cmd
}
