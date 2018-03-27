package camel

import "fmt"

// ==========================
//
// ProcessDefinition
//
//    WORK IN PROGRESS
//
// ==========================

// Process --
func (definition *RouteDefinition) Process() *ProcessDefinition {
	process := ProcessDefinition{}
	process.parent = definition

	definition.child = &process.RouteDefinition

	return &process
}

// ==========================
//
//
//
//
//
// ==========================

// ProcessDefinition --
type ProcessDefinition struct {
	RouteDefinition
}

// Fn --
func (definition *ProcessDefinition) Fn(consumer func(*Exchange)) *RouteDefinition {
	definition.parent.AddFactory(func(context *Context, parent Processor) (Processor, Service, error) {
		p := NewProcessorWithParent(parent, func(e *Exchange, out chan<- *Exchange) {
			consumer(e)
			out <- e
		})

		return p, nil, nil
	})

	return definition.parent
}

// Ref --
func (definition *ProcessDefinition) Ref(ref string) *RouteDefinition {
	definition.parent.AddFactory(func(context *Context, parent Processor) (Processor, Service, error) {
		registry := context.Registry()
		ifc, err := registry.Lookup(ref)

		if ifc != nil && err == nil {
			if consumer, ok := ifc.(func(*Exchange)); ok {
				p := NewProcessorWithParent(parent, func(e *Exchange, out chan<- *Exchange) {
					consumer(e)
					out <- e
				})

				return p, nil, nil
			}
		}

		// TODO: error handling
		return parent, nil, fmt.Errorf("Unsupported type for ref:%s, type=%T", ref, ifc)
	})

	return definition.parent
}
