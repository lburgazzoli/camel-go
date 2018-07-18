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
