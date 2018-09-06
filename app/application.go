// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"os"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/camel"
	"github.com/lburgazzoli/camel-go/logger"
	"github.com/lburgazzoli/camel-go/route"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
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

	logger.Log(zerolog.DebugLevel, "flow file is: %s", v.ConfigFileUsed())

	app := Application{}
	app.logger = logger.New("app")
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

// NewJS --
func NewJS(rd string) (*Application, error) {
	app := Application{}
	app.logger = logger.New("app")
	app.context = camel.NewContext()

	if routes, err := route.LoadFromJS(app.context, rd); err == nil {
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
	logger  zerolog.Logger
}

// Start --
func (app *Application) Start() {
	app.logger.Info().Msg("Start context")
	app.context.Start()
}

// Stop --
func (app *Application) Stop() {
	app.logger.Info().Msg("Stop context")
	app.context.Stop()
}
