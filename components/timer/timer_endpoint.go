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

// EndpointOptions --
type EndpointOptions struct {
	period time.Duration
}

// SetPeriod --
func (options *EndpointOptions) SetPeriod(period time.Duration) {
	options.period = period
}

// ==========================
//
// Functional Options
//
// ==========================

// EndpointOption --
type EndpointOption func(*EndpointOptions)

// Period --
func Period(value time.Duration) EndpointOption {
	return func(args *EndpointOptions) {
		args.period = value
	}
}

// ==========================
//
// Endpoint
//
// ==========================

func newEndpoint(component *Component, setters ...EndpointOption) (*timerEndpoint, error) {
	endpoint := timerEndpoint{}
	endpoint.component = component

	// Apply options
	for _, setter := range setters {
		setter(&endpoint.EndpointOptions)
	}

	return &endpoint, nil
}

type timerEndpoint struct {
	EndpointOptions

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
