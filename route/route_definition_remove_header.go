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

// RemoveHeader --
func (definition *RouteDefinition) RemoveHeader(key string) *RouteDefinition {
	d := RemoveHeadersDefinition{
		parent:     definition,
		children:   nil,
		headerKeys: []string{key},
	}

	definition.AddChild(&d)

	return definition
}

// RemoveHeaders --
func (definition *RouteDefinition) RemoveHeaders(keys []string) *RouteDefinition {
	d := RemoveHeadersDefinition{
		parent:     definition,
		children:   nil,
		headerKeys: keys,
	}

	definition.AddChild(&d)

	return definition
}

// ==========================
//
// FilterDefinition
//
// ==========================

// RemoveHeadersDefinition --
type RemoveHeadersDefinition struct {
	api.ContextAware
	ProcessingNode

	context  api.Context
	parent   *RouteDefinition
	children []Definition

	headerKeys []string
}

// SetContext --
func (definition *RemoveHeadersDefinition) SetContext(context api.Context) {
	definition.context = context
}

// Context --
func (definition *RemoveHeadersDefinition) Context() api.Context {
	return definition.context
}

// Parent --
func (definition *RemoveHeadersDefinition) Parent() Definition {
	return definition.parent
}

// Children --
func (definition *RemoveHeadersDefinition) Children() []Definition {
	return definition.children
}

// Processor ---
func (definition *RemoveHeadersDefinition) Processor() (api.Processor, error) {
	if definition.headerKeys != nil {
		p := processor.NewProcessingPipeline(func(exchange api.Exchange) {
			for _, k := range definition.headerKeys {
				exchange.Headers().Remove(k)
			}
		})

		return p, nil
	}

	return nil, nil
}
