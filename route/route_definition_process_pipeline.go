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

	"github.com/rs/zerolog"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/logger"
	"github.com/lburgazzoli/camel-go/processor"
)

// ==========================
//
// Extend RouteDefinition DSL
//
// ==========================

// Pipeline --
func (definition *RouteDefinition) Pipeline() *ProcessPipelineDefinition {
	d := ProcessPipelineDefinition{
		parent:     definition,
		children:   nil,
		processors: make([]api.ProcessingFnSupplier, 0),
	}

	definition.AddChild(&d)

	return &d
}

// ==========================
//
// ProcessPipelineDefinition
//
// ==========================

// ProcessPipelineDefinition --
type ProcessPipelineDefinition struct {
	api.ContextAware
	ProcessingNode

	parent   *RouteDefinition
	children []Definition

	context    api.Context
	processors []api.ProcessingFnSupplier
}

// SetContext --
func (definition *ProcessPipelineDefinition) SetContext(context api.Context) {
	definition.context = context
}

// Context --
func (definition *ProcessPipelineDefinition) Context() api.Context {
	return definition.context
}

// Parent --
func (definition *ProcessPipelineDefinition) Parent() Definition {
	return definition.parent
}

// Children --
func (definition *ProcessPipelineDefinition) Children() []Definition {
	return definition.children
}

// Processor ---
func (definition *ProcessPipelineDefinition) Processor() (api.Processor, error) {
	suppliers := make([]func(api.Exchange), len(definition.processors))

	for _, s := range definition.processors {
		p, e := s()

		if e != nil {
			logger.Log(zerolog.FatalLevel, e.Error())
		}

		if p != nil {
			suppliers = append(suppliers, p)
		}
	}

	if len(suppliers) == 1 {
		return processor.NewProcessingPipeline(suppliers[0]), nil
	} else if len(suppliers) > 1 {
		return processor.NewProcessingPipeline(suppliers[0], suppliers[1:]...), nil
	}

	return nil, fmt.Errorf("no element added to the pipeline")
}

// Fn --
func (definition *ProcessPipelineDefinition) Fn(processor func(api.Exchange)) *ProcessPipelineDefinition {
	// wrap processing function
	supplier := func() (func(api.Exchange), error) {
		return processor, nil
	}

	// add supplier
	definition.processors = append(definition.processors, supplier)

	return definition
}

// Ref --
func (definition *ProcessPipelineDefinition) Ref(ref string) *ProcessPipelineDefinition {
	// wrap processing function
	supplier := func() (func(api.Exchange), error) {
		registry := definition.context.Registry()
		ifc, found := registry.Lookup(ref)

		if ifc != nil && found {
			if p, ok := ifc.(func(api.Exchange)); ok {
				return p, nil
			}
		}

		return nil, fmt.Errorf("unable to resolve reference: %s", ref)
	}

	// add supplier
	definition.processors = append(definition.processors, supplier)

	return definition
}

// End --
func (definition *ProcessPipelineDefinition) End() *RouteDefinition {
	return definition.parent
}
