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

package log

import (
	"errors"
	"os"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/rs/zerolog"
)

// ==========================
//
// Options
//
// ==========================

// EndpointOptions --
type EndpointOptions struct {
	logger     string
	level      zerolog.Level
	logHeaders bool
}

// SetLogger --
func (options *EndpointOptions) SetLogger(logger string) {
	options.logger = logger
}

// SetLevel --
func (options *EndpointOptions) SetLevel(level zerolog.Level) {
	options.level = level
}

// SetLogHeaders --
func (options *EndpointOptions) SetLogHeaders(logHeaders bool) {
	options.logHeaders = logHeaders
}

// ==========================
//
// Functional Options
//
// ==========================

// EndpointOption --
type EndpointOption func(*EndpointOptions)

// Logger --
func Logger(value string) EndpointOption {
	return func(args *EndpointOptions) {
		args.logger = value
	}
}

// Level --
func Level(value zerolog.Level) EndpointOption {
	return func(args *EndpointOptions) {
		args.level = value
	}
}

// Headers --
func Headers(value bool) EndpointOption {
	return func(args *EndpointOptions) {
		args.logHeaders = value
	}
}

// ==========================
//
// Endpoint
//
// ==========================

func newEndpoint(component *Component, logger string, setters ...EndpointOption) (*logEndpoint, error) {
	endpoint := logEndpoint{}
	endpoint.component = component
	endpoint.logger = logger
	endpoint.level = zerolog.InfoLevel

	// Apply options
	for _, setter := range setters {
		setter(&endpoint.EndpointOptions)
	}

	return &endpoint, nil
}

type logEndpoint struct {
	EndpointOptions
	component *Component
}

func (endpoint *logEndpoint) Start() {
}

func (endpoint *logEndpoint) Stop() {
}

func (endpoint *logEndpoint) Stage() api.ServiceStage {
	return api.ServiceStageEndpoint
}

func (endpoint *logEndpoint) Component() api.Component {
	return endpoint.component
}

func (endpoint *logEndpoint) CreateProducer() (api.Producer, error) {
	// need to be replaced with better configuration from camel logging
	newlog := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger := newlog.With().Str("logger", endpoint.logger).Logger()

	return newLogProducer(endpoint, &logger), nil
}

func (endpoint *logEndpoint) CreateConsumer() (api.Consumer, error) {
	return nil, errors.New("log is Producer only")
}
