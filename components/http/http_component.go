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

package http

import (
	"fmt"
	ghttp "net/http"
	"net/url"

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
	Scheme string
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
		logger:         logger.New("http.Component"),
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

	if component.Scheme == "" {
		component.Scheme = "http"
	}
	url, err := url.Parse(fmt.Sprintf("%s://%s", component.Scheme, remaining))
	if err != nil {
		return nil, err
	}

	// Create the endpoint and set default values
	endpoint, err := newEndpoint(component, *url)
	if err != nil {
		return nil, err
	}

	// bind options to endpoint
	introspection.SetProperties(component.context, endpoint, options)

	// Check if method option was valid
	if endpoint.method != "" {
		if endpoint.method != ghttp.MethodGet &&
			endpoint.method != ghttp.MethodPost &&
			endpoint.method != ghttp.MethodPut &&
			endpoint.method != ghttp.MethodDelete &&
			endpoint.method != ghttp.MethodOptions &&
			endpoint.method != ghttp.MethodConnect &&
			endpoint.method != ghttp.MethodHead &&
			endpoint.method != ghttp.MethodPatch &&
			endpoint.method != ghttp.MethodTrace {
			return nil, fmt.Errorf("invalid HTTP method")
		}
	}

	return endpoint, nil
}

// ==========================
//
// Helpers
//
// ==========================

func (component *Component) doStart() {
	component.logger.Debug().Msg("Started")
}

func (component *Component) doStop() {
	component.logger.Debug().Msg("Stopped")
}
