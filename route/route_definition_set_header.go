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

// SetHeader --
func (definition *RouteDefinition) SetHeader(key string, val interface{}) *RouteDefinition {
	d := SetHeadersDefinition{
		parent:   definition,
		children: nil,
		headers:  map[string]interface{}{key: val},
	}

	definition.AddChild(&d)

	return definition
}

// SetHeaders --
func (definition *RouteDefinition) SetHeaders(headers map[string]interface{}) *RouteDefinition {
	d := SetHeadersDefinition{
		parent:   definition,
		children: nil,
		headers:  headers,
	}

	definition.AddChild(&d)

	return definition
}

// ==========================
//
// FilterDefinition
//
// ==========================

// SetHeadersDefinition --
type SetHeadersDefinition struct {
	api.ContextAware
	ProcessingNode

	context  api.Context
	parent   *RouteDefinition
	children []Definition

	headers map[string]interface{}
}

// SetContext --
func (definition *SetHeadersDefinition) SetContext(context api.Context) {
	definition.context = context
}

// Context --
func (definition *SetHeadersDefinition) Context() api.Context {
	return definition.context
}

// Parent --
func (definition *SetHeadersDefinition) Parent() Definition {
	return definition.parent
}

// Children --
func (definition *SetHeadersDefinition) Children() []Definition {
	return definition.children
}

// Processor ---
func (definition *SetHeadersDefinition) Processor() (api.Processor, error) {
	if definition.headers != nil {
		p := processor.NewProcessingPipeline(func(exchange api.Exchange) {
			for k, v := range definition.headers {
				exchange.Headers().Bind(k, v)
			}
		})

		return p, nil
	}

	return nil, nil
}
