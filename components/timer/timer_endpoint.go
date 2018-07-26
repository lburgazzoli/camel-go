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
	"errors"
	"time"

	"github.com/lburgazzoli/camel-go/api"
)

// ==========================
//
// Options
//
// ==========================

// Options --
type Options struct {
	period time.Duration
}

// SetPeriod --
func (options *Options) SetPeriod(period time.Duration) {
	options.period = period
}

// ==========================
//
// Functional Options
//
// ==========================

// Option --
type Option func(*Options)

// Period --
func Period(value time.Duration) Option {
	return func(args *Options) {
		args.period = value
	}
}

// ==========================
//
// Endpoint
//
// ==========================

func newEndpoint(component *Component, setters ...Option) (*timerEndpoint, error) {
	endpoint := timerEndpoint{}
	endpoint.component = component

	// Apply options
	for _, setter := range setters {
		setter(&endpoint.Options)
	}

	return &endpoint, nil
}

type timerEndpoint struct {
	Options

	component *Component
}

func (endpoint *timerEndpoint) Start() {
}

func (endpoint *timerEndpoint) Stop() {
}

func (endpoint *timerEndpoint) Stage() api.ServiceStage {
	return api.ServiceStageEndpoint
}

func (endpoint *timerEndpoint) Component() api.Component {
	return endpoint.component
}

func (endpoint *timerEndpoint) CreateProducer() (api.Producer, error) {
	return nil, errors.New("log is Consumer only")
}

func (endpoint *timerEndpoint) CreateConsumer() (api.Consumer, error) {
	return newTimerConsumer(endpoint), nil
}
