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
	"strconv"
	"time"

	ghttp "net/http"
	"net/url"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/introspection"
	"github.com/rs/zerolog"

	zlog "github.com/rs/zerolog/log"
)

// ==========================
//
//
//
// ==========================

// NewComponent --
func NewComponent() api.Component {
	component := &Component{
		logger:         zlog.With().Str("http", "http.Component").Logger(),
		serviceSupport: api.NewServiceSupport(),
	}

	component.serviceSupport.Transition(api.ServiceStatusSTOPPED, api.ServiceStatusSTARTED, component.doStart)
	component.serviceSupport.Transition(api.ServiceStatusSTARTED, api.ServiceStatusSTOPPED, component.doStop)

	return component
}

// ==========================
//
// Component
//
// ==========================

// Component --
type Component struct {
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
	var url *url.URL
	var err error

	if url, err = url.Parse("http://" + remaining); err != nil {
		return nil, err
	}

	// Create the endpoint and set default values
	endpoint := httpEndpoint{}
	endpoint.component = component
	endpoint.method = ghttp.MethodGet
	endpoint.connectionTimeout = 10 * time.Second
	endpoint.requestTimeout = 60 * time.Second
	endpoint.path = url.Path
	endpoint.port = 80

	if url.Hostname() != "" {
		endpoint.host = url.Hostname()
	}

	if url.Port() != "" {
		endpoint.port, err = strconv.Atoi(url.Port())

		if err != nil {
			return nil, err
		}
	}

	if endpoint.port == 443 && endpoint.scheme == "" {
		endpoint.scheme = "https"
	}

	if endpoint.scheme == "" {
		endpoint.scheme = "http"
	}

	// bind options to endpoint
	introspection.SetProperties(component.context, &endpoint, options)

	return &endpoint, nil
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
