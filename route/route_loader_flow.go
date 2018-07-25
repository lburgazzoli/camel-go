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

import (
	"fmt"
	"io/ioutil"

	"github.com/lburgazzoli/camel-go/api"
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
	context  api.Context
	flows    []Flow
	handlers map[string]StepHandler
}

// ==========================
//
// Initialization
//
// ==========================

// NewFlowLoader --
func NewFlowLoader(context api.Context, flows []Flow) *FlowLoader {
	loader := FlowLoader{
		context:  context,
		flows:    flows,
		handlers: make(map[string]StepHandler, 0),
	}

	loader.handlers["endpoint"] = EndpointStepHandler
	loader.handlers["process"] = ProcessStepHandler
	loader.handlers["filter"] = FilterStepHandler
	loader.handlers["header"] = SetHeaderStepHandler
	loader.handlers["headers"] = SetHeadersStepHandler

	return &loader
}

// Load --
func (loader *FlowLoader) Load() ([]*api.Route, error) {
	zlog.Info().Msgf("flows: %v", loader.flows)

	routes := make([]*api.Route, 0)

	for _, f := range loader.flows {
		var definition *RouteDefinition

		for i, s := range f.Steps {
			if i == 0 {
				definition = From(s["uri"].(string))
			} else {
				if t, ok := s["type"]; ok {
					h, e := findHandler(loader.handlers, t.(string))
					if e != nil {
						return nil, e
					}

					if r, e := h(s, definition); e == nil {
						definition = r
					} else {
						return nil, fmt.Errorf("Error handling step: %s, error: %v", s, e)
					}
				} else {
					return nil, fmt.Errorf("No Type defined for step: %v", s)
				}
			}
		}

		r, e := ToRoute(loader.context, definition)
		if e != nil {
			return nil, e
		}

		routes = append(routes, r)
	}

	return routes, nil
}

func (loader *FlowLoader) findHandler(stepType string) (StepHandler, error) {
	r := loader.context.Registry()
	h, f := r.Lookup(stepType)

	if f && h != nil {
		if s, ok := h.(StepHandler); ok {
			return s, nil
		}
	}

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
func LoadFlowFromYAMLFile(context api.Context, path string) ([]*api.Route, error) {
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

	return NewFlowLoader(context, integration.Flows).Load()
}

// LoadFlowFromViper --
func LoadFlowFromViper(context api.Context, v *viper.Viper) ([]*api.Route, error) {
	flows := make([]Flow, 0)
	err := v.UnmarshalKey("flows", &flows)

	if err != nil {
		return nil, err
	}

	return NewFlowLoader(context, flows).Load()
}
