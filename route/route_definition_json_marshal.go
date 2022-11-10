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
	"github.com/lburgazzoli/camel-go/logger"
	"github.com/lburgazzoli/camel-go/processor"
	"github.com/rs/zerolog"
)

// ==========================
//
// Extend RouteDefinition DSL
//
// ==========================

// JsonMarshal --
func (definition *RouteDefinition) JsonMarshal() *RouteDefinition {
	d := JsonMarshalDefinition{
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

// JsonMarshalDefinition --
type JsonMarshalDefinition struct {
	api.ContextAware
	ProcessingNode

	context  api.Context
	parent   *RouteDefinition
	children []Definition
}

// SetContext --
func (definition *JsonMarshalDefinition) SetContext(context api.Context) {
	definition.context = context
}

// Context --
func (definition *JsonMarshalDefinition) Context() api.Context {
	return definition.context
}

// Parent --
func (definition *JsonMarshalDefinition) Parent() Definition {
	return definition.parent
}

// Children --
func (definition *JsonMarshalDefinition) Children() []Definition {
	return definition.children
}

// Processor ---
func (definition *JsonMarshalDefinition) Processor() (api.Processor, error) {

	p := processor.NewProcessingPipeline(func(exchange api.Exchange) {

		body := exchange.Body()
		// convert body to json
		data, err := json.Marshal(body)
		if err != nil {
			logger.Log(zerolog.ErrorLevel, err.Error())
			exchange.SetError(err)
			return
		}
		// set the body to the json
		exchange.SetBody(string(data))
	})

	return p, nil

}
