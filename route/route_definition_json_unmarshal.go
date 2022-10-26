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
	"encoding/json"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/camel"
	"github.com/lburgazzoli/camel-go/logger"
	"github.com/lburgazzoli/camel-go/processor"
	"github.com/rs/zerolog"
)

// ==========================
//
// Extend RouteDefinition DSL
//
// ==========================

// JsonUnmarshal --
func (definition *RouteDefinition) JsonUnmarshal() *RouteDefinition {
	d := JsonUnmarshalDefinition{
		parent:   definition,
		children: nil,
	}

	definition.AddChild(&d)

	return definition
}

// ==========================
//
// FilterDefinition
//
// ==========================

// JsonUnmarshalDefinition --
type JsonUnmarshalDefinition struct {
	api.ContextAware
	ProcessingNode

	context  api.Context
	parent   *RouteDefinition
	children []Definition
}

// SetContext --
func (definition *JsonUnmarshalDefinition) SetContext(context api.Context) {
	definition.context = context
}

// Context --
func (definition *JsonUnmarshalDefinition) Context() api.Context {
	return definition.context
}

// Parent --
func (definition *JsonUnmarshalDefinition) Parent() Definition {
	return definition.parent
}

// Children --
func (definition *JsonUnmarshalDefinition) Children() []Definition {
	return definition.children
}

//TODO: error handling
// Processor ---
func (definition *JsonUnmarshalDefinition) Processor() (api.Processor, error) {

	p := processor.NewProcessingPipeline(func(exchange api.Exchange) {

		body := exchange.Body()
		tc := definition.context.TypeConverter()
		var data map[string]any

		// convert body to string
		str, err := tc(body, camel.TypeString)
		if err != nil {
			// do nothing here for the moment, we should fail the exchange
			logger.Log(zerolog.ErrorLevel, err.Error())
			return
		}

		// convert body string to map[string]any
		if err := json.Unmarshal([]byte(str.(string)), &data); err != nil {
			// do nothing here for the moment, we should fail the exchange
			logger.Log(zerolog.ErrorLevel, err.Error())
			return
		}

		// set the body to the map
		exchange.SetBody(data)
	})

	return p, nil

}
