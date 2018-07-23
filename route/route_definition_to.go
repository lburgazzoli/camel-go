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
	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/processor"
)

// ==========================
//
// Extend RouteDefinition DSL
//
// ==========================

// To --
func (definition *RouteDefinition) To(uri string) *RouteDefinition {
	d := ToDefinition{
		parent:   definition,
		children: nil,
		URI:      uri,
	}

	definition.AddChild(&d)

	return definition
}

// ==========================
//
// ToDefinition
//
// ==========================

// ToDefinition --
type ToDefinition struct {
	api.ContextAware
	ProcessingNode

	context  api.Context
	parent   *RouteDefinition
	children []Definition

	URI string
}

// SetContext --
func (definition *ToDefinition) SetContext(context api.Context) {
	definition.context = context
}

// Context --
func (definition *ToDefinition) Context() api.Context {
	return definition.context
}

// Parent --
func (definition *ToDefinition) Parent() Definition {
	return definition.parent
}

// Children --
func (definition *ToDefinition) Children() []Definition {
	return definition.children
}

// Processor ---
func (definition *ToDefinition) Processor() (api.Processor, error) {
	var err error
	var producer api.Producer
	var endpoint api.Endpoint

	if endpoint, err = api.NewEndpointFromURI(definition.context, definition.URI); err != nil {
		return nil, err
	}

	if producer, err = endpoint.CreateProducer(); err != nil {
		return nil, err
	}

	// TODO: re-engine
	endpoint.Start()

	return processor.NewProcessingService(producer, producer.Processor()), nil
}
