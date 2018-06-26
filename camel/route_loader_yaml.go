package camel

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	zlog "github.com/rs/zerolog/log"
)

// ==========================
//
// YAMLLoader
//
// ==========================

// YAMLLoader --
type YAMLLoader struct {
	data     []byte
	handlers map[string]StepHandler
}

// ==========================
//
// Initialization
//
// ==========================

// NewYAMLLoader --
func NewYAMLLoader(data []byte) RouteLoader {
	loader := YAMLLoader{
		data:     data,
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

func (loader *YAMLLoader) findHandler(stepType string) (StepHandler, error) {
	if h, ok := loader.handlers[stepType]; ok {
		return h, nil
	}

	return nil, fmt.Errorf("No StepHandler defined for type: %s", stepType)
}

// Load --
func (loader *YAMLLoader) Load() ([]Definition, error) {
	integration := Integration{}
	err := yaml.Unmarshal([]byte(loader.data), &integration)

	if err != nil {
		return nil, err
	}

	definitions := make([]Definition, 0)

	for _, f := range integration.Flows {
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

// ==========================
//
// Helpers
//
// ==========================

// LoadRouteFromYAMLFile --
func LoadRouteFromYAMLFile(path string) ([]Definition, error) {
	zlog.Debug().Msgf("Loading routes from:  %s", path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	loader := NewYAMLLoader(data)
	return loader.Load()
}
