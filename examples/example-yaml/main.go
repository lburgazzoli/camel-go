package main

import (
	"bufio"
	"os"

	"github.com/lburgazzoli/camel-go/camel"
	_ "github.com/lburgazzoli/camel-go/components/log"
	_ "github.com/lburgazzoli/camel-go/components/timer"

	"github.com/spf13/viper"

	zlog "github.com/rs/zerolog/log"
)

// ==========================
//
// Main
//
// ==========================

// ConfigureViper --
func ConfigureViper() *viper.Viper {
	v := viper.New()
	v.SetConfigName("configuration")
	v.SetConfigType("yaml")

	if wd, err := os.Getwd(); err == nil {
		v.AddConfigPath(wd)
	}

	if err := v.ReadInConfig(); err != nil {
		zlog.Panic().Msgf("fatal error config file: %s", err)
	}

	return v
}

func main() {
	context := camel.NewContext()
	v := ConfigureViper()

	for _, p := range v.GetStringSlice("plugins") {
		context.Registry().AddLoader(camel.NewPluginRegistryLoader(p))
	}

	if routes, err := camel.LoadFlowFromViper(context, v); err == nil {
		for _, r := range routes {
			context.AddRoute(r)
		}

		zlog.Info().Msg("Start context")
		context.Start()

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		zlog.Info().Msg("Stop context")
		context.Stop()
	} else {
		zlog.Info().Msgf("%s", err)
	}
}
