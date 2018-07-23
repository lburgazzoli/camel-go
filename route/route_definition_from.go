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

// FromDefinition --
type FromDefinition struct {
	api.ContextAware
	ProcessingNode

	context  api.Context
	parent   Definition
	children []Definition

	URI string
}

// SetContext --
func (definition *FromDefinition) SetContext(context api.Context) {
	definition.context = context
}

// Context --
func (definition *FromDefinition) Context() api.Context {
	return definition.context
}

// Parent --
func (definition *FromDefinition) Parent() Definition {
	return definition.parent
}

// Children --
func (definition *FromDefinition) Children() []Definition {
	return definition.children
}

// Processor ---
func (definition *FromDefinition) Processor() (api.Processor, error) {
	var err error
	var consumer api.Consumer
	var endpoint api.Endpoint

	if endpoint, err = api.NewEndpointFromURI(definition.context, definition.URI); err != nil {
		return nil, nil
	}

	// TODO: re-engine
	endpoint.Start()

	if consumer, err = endpoint.CreateConsumer(); err != nil {
		return nil, nil
	}

	return processor.NewProcessingService(consumer, consumer.Processor()), nil
}
