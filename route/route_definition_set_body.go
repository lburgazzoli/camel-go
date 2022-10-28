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
	"github.com/lburgazzoli/camel-go/logger"
	"github.com/lburgazzoli/camel-go/processor"
	"github.com/rs/zerolog"
)

// ==========================
//
// Extend RouteDefinition DSL
//
// ==========================

// SetBody --
func (definition *RouteDefinition) SetBody(body any) *RouteDefinition {
	d := SetBodyDefinition{
		parent:   definition,
		children: nil,
		body:     body,
	}

	definition.AddChild(&d)

	return definition
}

// ==========================
//
// FilterDefinition
//
// ==========================

// SetBodyDefinition --
type SetBodyDefinition struct {
	api.ContextAware
	ProcessingNode

	context  api.Context
	parent   *RouteDefinition
	children []Definition

	body any
}

// SetContext --
func (definition *SetBodyDefinition) SetContext(context api.Context) {
	definition.context = context
}

// Context --
func (definition *SetBodyDefinition) Context() api.Context {
	return definition.context
}

// Parent --
func (definition *SetBodyDefinition) Parent() Definition {
	return definition.parent
}

// Children --
func (definition *SetBodyDefinition) Children() []Definition {
	return definition.children
}

//TODO: error handling
// Processor ---
func (definition *SetBodyDefinition) Processor() (api.Processor, error) {
	if definition.body != nil {
		p := processor.NewProcessingPipeline(func(exchange api.Exchange) {

			// Check if the body is an expression
			if expr, ok := definition.body.(api.Expression); ok {

				body, err := expr.Evaluate(exchange)

				if err != nil {
					// do nothing here for the moment, we should fail the exchange
					logger.Log(zerolog.ErrorLevel, err.Error())
					body = expr.Raw()
				}
				exchange.SetBody(body)
			} else {
				exchange.SetBody(definition.body)
			}
		})

		return p, nil
	}

	return nil, nil
}
