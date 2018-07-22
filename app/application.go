package app

import (
	"os"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/camel"
	"github.com/lburgazzoli/camel-go/route"
	"github.com/spf13/viper"

	zlog "github.com/rs/zerolog/log"
)

// New --
func New(config string) (*Application, error) {
	v := viper.New()

	if config != "" {
		v.SetConfigFile(config)
	} else {
		v.SetConfigName("flow")
		v.SetConfigType("yaml")

		if wd, err := os.Getwd(); err == nil {
			v.AddConfigPath(wd)
		}
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	zlog.Debug().Msgf("flow file is: %s", v.ConfigFileUsed())

	app := Application{}
	app.context = camel.NewContext()

	for _, p := range v.GetStringSlice("plugins") {
		app.context.Registry().AddLoader(camel.NewPluginRegistryLoader(p))
	}

	if routes, err := route.LoadFlowFromViper(app.context, v); err == nil {
		for _, r := range routes {
			app.context.AddRoute(r)
		}
	} else {
		return nil, err
	}

	return &app, nil
}

// Application --
type Application struct {
	context api.Context
}

// Start --
func (app *Application) Start() {
	zlog.Info().Msg("Start context")
	app.context.Start()
}

// Stop --
func (app *Application) Stop() {
	zlog.Info().Msg("Stop context")
	app.context.Stop()
}
