package camel

import (
	"fmt"

	zlog "github.com/rs/zerolog/log"
)

// ==========================
//
//
//
// ==========================

// FlowLoader --
type FlowLoader struct {
	handlers map[string]StepHandler
}

// ==========================
//
// Initialization
//
// ==========================

// NewFlowLoader --
func NewFlowLoader() FlowLoader {
	loader := FlowLoader{
		handlers: make(map[string]StepHandler, 0),
	}

	loader.handlers["endpoint"] = EndpointStepHandler
	loader.handlers["process"] = ProcessStepHandler
	loader.handlers["filter"] = FilterStepHandler

	return loader
}

// ToDefinition --
func (loader *FlowLoader) definition(flows []Flow) ([]Definition, error) {
	zlog.Info().Msgf("flows: %v", flows)

	definitions := make([]Definition, 0)

	for _, f := range flows {
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
