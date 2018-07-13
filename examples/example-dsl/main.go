package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"

	"github.com/lburgazzoli/camel-go/camel"
	"github.com/lburgazzoli/camel-go/components/log"
	"github.com/lburgazzoli/camel-go/components/timer"
	"github.com/lburgazzoli/camel-go/types"
	"github.com/spf13/viper"

	zlog "github.com/rs/zerolog/log"
)

// ==========================
//
// Functions
//
// ==========================

func simpleProcess(e *camel.Exchange) {
	e.SetHeader("ref.header", rand.Int())
}

func simpleFilter(e *camel.Exchange) bool {
	count := e.HeaderAs("timer.fire.count", types.TypeInt).(int)
	return count != 4
}

func simpleProcessorFn(e *camel.Exchange) {
	e.SetBody(fmt.Sprintf("random body: %d", rand.Int()))
}

func simpleFilterFn(e *camel.Exchange) bool {
	count := e.HeaderAs("timer.fire.count", types.TypeInt).(int)
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

	def := camel.From("timer:start?period=1s").
		Process().Fn(simpleProcessorFn).
		Process().Ref("refProcessor").
		Filter().Fn(simpleFilterFn).
		Filter().Ref("refFilter").
		To("log:test?logHeaders=true")

	context.AddRouteDefinition(def)

	zlog.Info().Msg("Start context")
	context.Start()

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	zlog.Info().Msg("Stop context")
	context.Stop()
}
