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

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/camel"
	"github.com/lburgazzoli/camel-go/components/log"
	"github.com/lburgazzoli/camel-go/components/timer"
	"github.com/lburgazzoli/camel-go/route"
	"github.com/spf13/viper"

	zlog "github.com/rs/zerolog/log"
)

// ==========================
//
// Functions
//
// ==========================

func simpleProcess(e api.Exchange) {
	e.Headers().Bind("ref.header", rand.Int())
}

func simpleFilter(e api.Exchange) bool {
	c, ok := e.Headers().LookupAs("timer.fire.count", camel.TypeInt)
	if !ok {
		zlog.Panic().Msg("Unable to convert header")
	}

	count, ok := c.(int)
	if !ok {
		zlog.Panic().Msg("Unable to convert header")
	}

	return count != 4
}

func simpleProcessorFn(e api.Exchange) {
	e.SetBody(fmt.Sprintf("random body: %d", rand.Int()))
}

func simpleFilterFn(e api.Exchange) bool {
	c, ok := e.Headers().LookupAs("timer.fire.count", camel.TypeInt)
	if !ok {
		zlog.Panic().Msg("Unable to convert header")
	}

	count, ok := c.(int)
	if !ok {
		zlog.Panic().Msg("Unable to convert header")
	}

	return count%2 == 0
}

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

	context.Registry().Bind("log", log.NewComponent())
	context.Registry().Bind("timer", timer.NewComponent())
	context.Registry().Bind("log", log.NewComponent())
	context.Registry().Bind("timer", timer.NewComponent())
	context.Registry().Bind("refProcessor", simpleProcess)
	context.Registry().Bind("refFilter", simpleFilter)

	def := route.From("timer:start?period=1s").
		Process().Fn(simpleProcessorFn).
		Process().Ref("refProcessor").
		Filter().Fn(simpleFilterFn).
		Filter().Ref("refFilter").
		To("log:test?logHeaders=true")

	r, err := route.ToRoute(context, def)
	if err != nil {
		zlog.Panic().Msg("Unable to load route")
	}

	context.AddRoute(r)

	zlog.Info().Msg("Start context")
	context.Start()

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	zlog.Info().Msg("Stop context")
	context.Stop()
}
