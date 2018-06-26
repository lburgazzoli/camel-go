package camel

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	zlog "github.com/rs/zerolog/log"
)

// ==========================
//
// Model
//
// ==========================

// Step --
type Step map[string]interface{}

// StepHandler --
type StepHandler func(step Step, route *RouteDefinition) (*RouteDefinition, error)

// Flow --
type Flow struct {
	Steps []Step `yaml:"steps"`
}

// Integration --
type Integration struct {
	ID    string `yaml:"id"`
	Name  string `yaml:"name"`
	Flows []Flow `yaml:"flows"`
}

// TypedStep --
type TypedStep struct {
	Type string `yaml:"type"`
}

// ==========================
//
// YAMLLoader
//
// ==========================

// YAMLLoader --
type YAMLLoader struct {
	handlers map[string]StepHandler
}

// ==========================
//
// Initialization
//
// ==========================

// NewYAMLLoader --
func NewYAMLLoader() *YAMLLoader {
	loader := YAMLLoader{
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
func (loader *YAMLLoader) Load(data []byte) ([]Definition, error) {
	integration := Integration{}
	err := yaml.Unmarshal([]byte(data), &integration)

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

	loader := NewYAMLLoader()
	return loader.Load(data)
}
