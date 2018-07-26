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

// Options --
type Options struct {
	logger     string
	level      zerolog.Level
	logHeaders bool
}

// SetLogger --
func (options *Options) SetLogger(logger string) {
	options.logger = logger
}

// SetLevel --
func (options *Options) SetLevel(level zerolog.Level) {
	options.level = level
}

// SetLogHeaders --
func (options *Options) SetLogHeaders(logHeaders bool) {
	options.logHeaders = logHeaders
}

// ==========================
//
// Functional Options
//
// ==========================

// Option --
type Option func(*Options)

// Logger --
func Logger(value string) Option {
	return func(args *Options) {
		args.logger = value
	}
}

// Level --
func Level(value zerolog.Level) Option {
	return func(args *Options) {
		args.level = value
	}
}

// Headers --
func Headers(value bool) Option {
	return func(args *Options) {
		args.logHeaders = value
	}
}

// ==========================
//
// Endpoint
//
// ==========================

func newEndpoint(component *Component, logger string, setters ...Option) (*logEndpoint, error) {
	endpoint := logEndpoint{}
	endpoint.component = component
	endpoint.logger = logger
	endpoint.level = zerolog.InfoLevel

	// Apply options
	for _, setter := range setters {
		setter(&endpoint.Options)
	}

	return &endpoint, nil
}

type logEndpoint struct {
	Options
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
