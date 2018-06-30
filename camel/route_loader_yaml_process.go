package camel

import (
	"github.com/lburgazzoli/camel-go/module"
	"github.com/mitchellh/mapstructure"
	zlog "github.com/rs/zerolog/log"
)

// ProcessStep --
type ProcessStep struct {
	TypedStep

	Ref      string `yaml:"ref"`
	Location string `yaml:"location"`
}

// ProcessStepHandler --
func ProcessStepHandler(step Step, route *RouteDefinition) (*RouteDefinition, error) {
	var impl ProcessStep

	// not really needed, added for testing purpose
	err := mapstructure.Decode(step, &impl)
	if err != nil {
		return nil, err
	}

	zlog.Debug().Msgf("handle process: step=<%v>, impl=<%+v>", step, impl)

	if impl.Location != "" {
		symbol, err := module.LoadSymbol(impl.Location, impl.Ref)
		if err != nil {
			return nil, err
		}

		return route.Process().Fn(symbol.(func(*Exchange))), nil
	}

	return route.Process().Ref(impl.Ref), nil
}
