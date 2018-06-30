package camel

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/viper"

	zlog "github.com/rs/zerolog/log"
	yaml "gopkg.in/yaml.v1"
)

// ==========================
//
//
//
// ==========================

// FlowLoader --
type FlowLoader struct {
	flows    []Flow
	handlers map[string]StepHandler
}

// ==========================
//
// Initialization
//
// ==========================

// NewFlowLoader --
func NewFlowLoader(flows []Flow) *FlowLoader {
	loader := FlowLoader{
		flows:    flows,
		handlers: make(map[string]StepHandler, 0),
	}

	loader.handlers["endpoint"] = EndpointStepHandler
	loader.handlers["process"] = ProcessStepHandler
	loader.handlers["filter"] = FilterStepHandler

	return &loader
}

// Load --
func (loader *FlowLoader) Load() ([]Definition, error) {
	zlog.Info().Msgf("flows: %v", loader.flows)

	definitions := make([]Definition, 0)

	for _, f := range loader.flows {
		var route *RouteDefinition

		for i, s := range f.Steps {
			if i == 0 {
				route = From(s["uri"].(string))
			} else {
				if t, ok := s["type"]; ok {
					h, e := findHandler(loader.handlers, t.(string))
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

func (loader *FlowLoader) findHandler(stepType string) (StepHandler, error) {
	if h, ok := loader.handlers[stepType]; ok {
		return h, nil
	}

	return nil, fmt.Errorf("No StepHandler defined for type: %s", stepType)
}

// ==========================
//
// Helpers
//
// ==========================

// LoadFlowFromYAMLFile --
func LoadFlowFromYAMLFile(path string) ([]Definition, error) {
	zlog.Debug().Msgf("Loading routes from:  %s", path)

	var err error
	var data []byte

	data, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	integration := Integration{}
	err = yaml.Unmarshal([]byte(data), &integration)

	if err != nil {
		return nil, err
	}

	return NewFlowLoader(integration.Flows).Load()
}

// LoadFlowFromViper --
func LoadFlowFromViper(v *viper.Viper) ([]Definition, error) {
	flows := make([]Flow, 0)
	err := v.UnmarshalKey("flows", &flows)

	if err != nil {
		return nil, err
	}

	return NewFlowLoader(flows).Load()
}
