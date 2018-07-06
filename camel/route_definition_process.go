package camel

import (
	"fmt"

	"github.com/lburgazzoli/camel-go/api"
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
	parent   *RouteDefinition
	children []Definition

	processor    func(api.Exchange)
	processorRef string
}

// Parent --
func (definition *ProcessDefinition) Parent() Definition {
	return definition.parent
}

// Children --
func (definition *ProcessDefinition) Children() []Definition {
	return definition.children
}

// Unwrap ---
func (definition *ProcessDefinition) Unwrap(context api.Context, parent api.Processor) (api.Processor, api.Service, error) {
	if definition.processor != nil {
		p := api.NewProcessorWithParent(parent, func(e api.Exchange, out chan<- api.Exchange) {
			definition.processor(e)

			out <- e
		})

		return p, nil, nil
	}

	if definition.processorRef != "" {
		registry := context.Registry()
		ifc, found := registry.Lookup(definition.processorRef)

		if ifc != nil && found {
			if processor, ok := ifc.(func(e api.Exchange)); ok {
				p := api.NewProcessorWithParent(parent, func(e api.Exchange, out chan<- api.Exchange) {
					processor(e)

					out <- e
				})

				return p, nil, nil
			}
		}

		var err error

		if !found {
			err = fmt.Errorf("Unsupported type for ref:%s, type=%T", definition.processorRef, ifc)
		}

		// TODO: error handling
		return nil, nil, err
	}

	return nil, nil, nil

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
