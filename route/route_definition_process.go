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
	"fmt"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/processor"
)

// ==========================
//
// Extend RouteDefinition DSL
//
// ==========================

// Process --
func (definition *RouteDefinition) Process() *ProcessDefinition {
	d := ProcessDefinition{
		parent:   definition,
		children: nil,
	}

	definition.AddChild(&d)

	return &d
}

// ==========================
//
// ProcessDefinition
//
// ==========================

// ProcessDefinition --
type ProcessDefinition struct {
	api.ContextAware
	ProcessingNode

	parent   *RouteDefinition
	children []Definition

	context      api.Context
	processor    func(api.Exchange)
	processorRef string
}

// SetContext --
func (definition *ProcessDefinition) SetContext(context api.Context) {
	definition.context = context
}

// Context --
func (definition *ProcessDefinition) Context() api.Context {
	return definition.context
}

// Parent --
func (definition *ProcessDefinition) Parent() Definition {
	return definition.parent
}

// Children --
func (definition *ProcessDefinition) Children() []Definition {
	return definition.children
}

// Processor ---
func (definition *ProcessDefinition) Processor() (api.Processor, error) {
	if definition.processor != nil {
		return processor.NewProcessingPipeline(definition.processor), nil
	}

	if definition.processorRef != "" {
		registry := definition.context.Registry()
		ifc, found := registry.Lookup(definition.processorRef)

		if ifc != nil && found {
			if p, ok := ifc.(func(e api.Exchange)); ok {
				return processor.NewProcessingPipeline(p), nil
			}
		}

		var err error

		if !found {
			err = fmt.Errorf("Unsupported type for ref:%s, type=%T", definition.processorRef, ifc)
		}

		// TODO: error handling
		return nil, err
	}

	return nil, nil
}

// Fn --
func (definition *ProcessDefinition) Fn(processor func(api.Exchange)) *RouteDefinition {
	definition.processor = processor
	return definition.parent
}

// Ref --
func (definition *ProcessDefinition) Ref(ref string) *RouteDefinition {
	definition.processorRef = ref
	return definition.parent
}
