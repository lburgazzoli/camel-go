package camel

import (
	"fmt"
	"os"
	"plugin"

	"github.com/mitchellh/mapstructure"
	zlog "github.com/rs/zerolog/log"
)

// FilterStep --
type FilterStep struct {
	TypedStep

	Ref      string `yaml:"ref"`
	Location string `yaml:"location"`
}

// FilterStepHandler --
func FilterStepHandler(step Step, route *RouteDefinition) (*RouteDefinition, error) {
	var impl FilterStep

	// not really needed, added for testing purpose
	err := mapstructure.Decode(step, &impl)
	if err != nil {
		return nil, err
	}

	zlog.Debug().Msgf("handle filter: step=<%v>, impl=<%+v>", step, impl)

	if impl.Location != "" {
		_, err := os.Stat(impl.Location)

		if os.IsNotExist(err) {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		plug, err := plugin.Open(impl.Location)
		if err != nil {
			return nil, fmt.Errorf("failed to open plugin %s: %v", impl.Location, err)
		}

		symbol, err := plug.Lookup(impl.Ref)
		if err != nil {
			return nil, fmt.Errorf("plugin %s does not export symbol \"%s\"", impl.Location, impl.Ref)
		}

		return route.Filter().Fn(symbol.(func(*Exchange) bool)), nil
	}

	return route.Filter().Ref(impl.Ref), nil
}
