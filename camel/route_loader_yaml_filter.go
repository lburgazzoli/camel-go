package camel

import (
	"github.com/mitchellh/mapstructure"
	zlog "github.com/rs/zerolog/log"
)

// FilterStep --
type FilterStep struct {
	TypedStep

	Ref string `yaml:"ref"`
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
	return route.Filter().Ref(impl.Ref), nil
}
