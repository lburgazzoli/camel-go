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

package timer

import (
	"fmt"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/introspection"
	"github.com/lburgazzoli/camel-go/logger"

	"github.com/rs/zerolog"
)

// ==========================
//
// Options
//
// ==========================

// ComponentOptions --
type ComponentOptions struct {
}

// ==========================
//
// Functional Options
//
// ==========================

// ComponentOption --
type ComponentOption func(*ComponentOptions)

// ==========================
//
// Component
//
// ==========================

// NewComponent --
func NewComponent(setters ...ComponentOption) api.Component {
	component := Component{
		logger:         logger.New("timer.Component"),
		serviceSupport: api.NewServiceSupport(),
	}

	// Apply options
	for _, setter := range setters {
		setter(&component.ComponentOptions)
	}

	component.serviceSupport.Transition(api.ServiceStatusSTOPPED, api.ServiceStatusSTARTED, component.doStart)
	component.serviceSupport.Transition(api.ServiceStatusSTARTED, api.ServiceStatusSTOPPED, component.doStop)

	return &component
}

// Component --
type Component struct {
	ComponentOptions

	logger         zerolog.Logger
	serviceSupport api.ServiceSupport
	context        api.Context
}

// SetContext --
func (component *Component) SetContext(context api.Context) {
	component.context = context
}

// Context --
func (component *Component) Context() api.Context {
	return component.context
}

// Start --
func (component *Component) Start() {
	component.serviceSupport.To(api.ServiceStatusSTARTED)
}

// Stop --
func (component *Component) Stop() {
	component.serviceSupport.To(api.ServiceStatusSTOPPED)
}

// Stage --
func (component *Component) Stage() api.ServiceStage {
	return api.ServiceStageComponent
}

// CreateEndpoint --
func (component *Component) CreateEndpoint(remaining string, options map[string]interface{}) (api.Endpoint, error) {
	// endpoint option validation
	if _, ok := options["period"]; !ok {
		return nil, fmt.Errorf("missing mandatory option: period")
	}

	// Create the endpoint and set default values
	endpoint, err := newEndpoint(component)
	if err != nil {
		return nil, err
	}

	// bind options to endpoint
	introspection.SetProperties(component.context, endpoint, options)

	return endpoint, nil
}

// ==========================
//
// Helpers
//
// ==========================

func (component *Component) doStart() {
	component.logger.Info().Msg("Started")
}

func (component *Component) doStop() {
	component.logger.Info().Msg("Stopped")
}
