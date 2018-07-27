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
	"fmt"

	"github.com/rs/zerolog"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/camel"
	"github.com/lburgazzoli/camel-go/logger"
	"github.com/lburgazzoli/camel-go/processor"
	"github.com/oliveagle/jsonpath"
)

// ==========================
//
// Extend RouteDefinition DSL
//
// ==========================

// Filter --
func (definition *RouteDefinition) Filter() *FilterDefinition {
	d := FilterDefinition{
		parent:   definition,
		children: nil,
	}

	definition.AddChild(&d)

	return &d
}

// ==========================
//
// FilterDefinition
//
// ==========================

// FilterDefinition --
type FilterDefinition struct {
	api.ContextAware
	ProcessingNode

	context  api.Context
	parent   *RouteDefinition
	children []Definition

	predicate    func(api.Exchange) bool
	predicateRef string
}

// SetContext --
func (definition *FilterDefinition) SetContext(context api.Context) {
	definition.context = context
}

// Context --
func (definition *FilterDefinition) Context() api.Context {
	return definition.context
}

// Parent --
func (definition *FilterDefinition) Parent() Definition {
	return definition.parent
}

// Children --
func (definition *FilterDefinition) Children() []Definition {
	return definition.children
}

// Processor ---
func (definition *FilterDefinition) Processor() (api.Processor, error) {
	if definition.predicate != nil {
		return processor.NewFilteringPipeline(definition.predicate), nil
	}

	if definition.predicateRef != "" {
		registry := definition.context.Registry()
		ifc, found := registry.Lookup(definition.predicateRef)

		if ifc != nil && found {
			if p, ok := ifc.(func(e api.Exchange) bool); ok {
				return processor.NewFilteringPipeline(p), nil
			}
		}

		var err error

		if !found {
			err = fmt.Errorf("Unsupported type for ref:%s, type=%T", definition.predicateRef, ifc)
		}

		// TODO: error handling
		return nil, err
	}

	return nil, nil
}

// Fn --
func (definition *FilterDefinition) Fn(predicate func(api.Exchange) bool) *RouteDefinition {
	definition.predicate = predicate
	return definition.parent
}

// Ref --
func (definition *FilterDefinition) Ref(ref string) *RouteDefinition {
	definition.predicateRef = ref
	return definition.parent
}

// JSONPath --
func (definition *FilterDefinition) JSONPath(expression string) *RouteDefinition {
	path, err := jsonpath.Compile(expression)
	if err != nil {
		logger.Log(zerolog.FatalLevel, "unable to compile expression: %s", expression)
	}

	definition.predicate = func(e api.Exchange) bool {
		if b := e.BodyAs(camel.TypeString); b != nil {
			var body string
			var data interface{}

			// if we are here, checking the type conversion
			// is not needed as BodyAs would fail if the
			// conversion is not possible
			body = b.(string)

			// need to define a type and relate conversion
			// for arrays
			json.Unmarshal([]byte(body), &data)

			res, err := path.Lookup(data)
			if err != nil {
				return false
			}

			//
			return res != nil
		}

		return false
	}

	return definition.parent
}
