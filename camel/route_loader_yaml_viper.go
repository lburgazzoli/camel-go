package camel

import (
	"fmt"

	"github.com/spf13/viper"

	zlog "github.com/rs/zerolog/log"
)

// ==========================
//
// ViperYAMLLoader
//
// ==========================

// ViperYAMLLoader --
type ViperYAMLLoader struct {
	viper    *viper.Viper
	handlers map[string]StepHandler
}

// ==========================
//
// Initialization
//
// ==========================

// NewViperYAMLLoader --
func NewViperYAMLLoader(viper *viper.Viper) RouteLoader {
	loader := ViperYAMLLoader{
		viper:    viper,
		handlers: make(map[string]StepHandler, 0),
	}

	loader.handlers["endpoint"] = EndpointStepHandler
	loader.handlers["process"] = ProcessStepHandler
	loader.handlers["filter"] = FilterStepHandler

	return &loader
}

// ==========================
//
// Implementation
//
// ==========================

func (loader *ViperYAMLLoader) findHandler(stepType string) (StepHandler, error) {
	if h, ok := loader.handlers[stepType]; ok {
		return h, nil
	}

	return nil, fmt.Errorf("No StepHandler defined for type: %s", stepType)
}

// Load --
func (loader *ViperYAMLLoader) Load() ([]Definition, error) {
	flows := make([]Flow, 0)
	err := loader.viper.UnmarshalKey("flows", &flows)

	zlog.Info().Msgf("flows: %v", flows)
	if err != nil {
		return nil, err
	}

	definitions := make([]Definition, 0)

	for _, f := range flows {
		var route *RouteDefinition

		for i, s := range f.Steps {
			if i == 0 {
				route = From(s["uri"].(string))
			} else {
				if t, ok := s["type"]; ok {
					h, e := loader.findHandler(t.(string))
					if e != nil {
						return nil, e
					}

					if r, e := h(s, route); e == nil {
						route = r
					} else {
						return nil, fmt.Errorf("Error handling step: %s, error: %v", s, e)
					}
				} else {
					return nil, fmt.Errorf("No Type defined for step: %v", s)
				}
			}
		}

		definitions = append(definitions, route)
	}

	return definitions, nil
}
