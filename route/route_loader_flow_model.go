package route

import "fmt"

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
	ID    string `yaml:"id"`
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}

// Integration --
type Integration struct {
	Flows []Flow `yaml:"flows"`
}

// TypedStep --
type TypedStep struct {
	Type string `yaml:"type"`
}

// ==========================
//
// Model
//
// ==========================

func findHandler(handlers map[string]StepHandler, stepType string) (StepHandler, error) {
	if h, ok := handlers[stepType]; ok {
		return h, nil
	}

	return nil, fmt.Errorf("No StepHandler defined for type: %s", stepType)
}

// ToDefinition --
func (flow *Flow) ToDefinition(handlers map[string]StepHandler) (Definition, error) {

	var route *RouteDefinition

	for i, s := range flow.Steps {
		if i == 0 {
			route = From(s["uri"].(string))
		} else {
			if t, ok := s["type"]; ok {
				h, e := findHandler(handlers, t.(string))
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

	return route, nil
}
