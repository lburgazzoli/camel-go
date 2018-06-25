package camel

import (
	"github.com/spf13/viper"
)

// LoadRoute --
func LoadRoute(v *viper.Viper) (Definition, error) {
	flow := Flow{}
	err := v.UnmarshalKey("flow", &flow)

	if err != nil {
		return nil, err
	}

	var route *RouteDefinition

	if err == nil {
		for i, s := range flow.Steps {
			if i == 0 {
				route = From(s.URI)
			} else {
				route = route.To(s.URI)
			}
		}
	}

	return route, nil
}
