package cmd

import (
	"os"
	"os/signal"

	"github.com/lburgazzoli/camel-go/camel"
	"github.com/lburgazzoli/camel-go/route"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	zlog "github.com/rs/zerolog/log"

	// import components
	_ "github.com/lburgazzoli/camel-go/components/log"
	_ "github.com/lburgazzoli/camel-go/components/timer"
)

var flow string

func init() {
	runCmd.Flags().StringVarP(&flow, "flow", "f", "", "flow to run")

	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run",
	Long:  `run`,
	Run: func(cmd *cobra.Command, args []string) {

		context := camel.NewContext()
		v := configureViper()

		for _, p := range v.GetStringSlice("plugins") {
			context.Registry().AddLoader(camel.NewPluginRegistryLoader(p))
		}

		if routes, err := route.LoadFlowFromViper(context, v); err == nil {
			for _, r := range routes {
				context.AddRoute(r)
			}

			zlog.Info().Msg("Start context")
			context.Start()

			waitForCtrlC()

			zlog.Info().Msg("Stop context")
			context.Stop()
		} else {
			zlog.Info().Msgf("%s", err)
		}
	},
}

// WaitForCtrlC --
func waitForCtrlC() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

// ConfigureViper --
func configureViper() *viper.Viper {
	v := viper.New()

	if flow != "" {
		zlog.Debug().Msgf("flow file is: %s", flow)
		v.SetConfigFile(flow)
	} else {
		v.SetConfigName("configuration")
		v.SetConfigType("yaml")

		if wd, err := os.Getwd(); err == nil {
			v.AddConfigPath(wd)
		}
	}

	if err := v.ReadInConfig(); err != nil {
		zlog.Panic().Msgf("fatal error config file: %s", err)
	}

	return v
}
