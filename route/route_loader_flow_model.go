// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
