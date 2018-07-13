package route

import (
	"github.com/mitchellh/mapstructure"
	zlog "github.com/rs/zerolog/log"
)

// EndpointStep --
type EndpointStep struct {
	TypedStep

	URI string `yaml:"uri"`
}

// EndpointStepHandler --
func EndpointStepHandler(step Step, route *RouteDefinition) (*RouteDefinition, error) {
	var impl EndpointStep

	// not really needed, added for testing purpose
	err := mapstructure.Decode(step, &impl)
	if err != nil {
		return nil, err
	}

	zlog.Debug().Msgf("handle endpoint: step=<%v>, impl=<%+v>", step, impl)
	return route.To(impl.URI), nil
}
