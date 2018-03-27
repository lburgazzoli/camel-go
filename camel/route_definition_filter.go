package camel

import (
	"fmt"
)

// ==========================
//
// FilterDefinition
//
//    WORK IN PROGRESS
//
// ==========================

// Filter --
func (definition *RouteDefinition) Filter() *FilterDefinition {
	filter := FilterDefinition{}
	filter.parent = definition

	definition.child = &filter.RouteDefinition

	return &filter
}

// ==========================
//
//
//
//
//
// ==========================

// FilterDefinition --
type FilterDefinition struct {
	RouteDefinition
}

// Fn --
func (definition *FilterDefinition) Fn(predicate func(*Exchange) bool) *RouteDefinition {
	definition.parent.AddFactory(func(context *Context, parent Processor) (Processor, Service, error) {
		fn := func(e *Exchange, out chan<- *Exchange) {
			if predicate(e) {
				out <- e
			}
		}

		p := NewProcessor(fn)
		p.Parent(parent)

		return p, nil, nil
	})

	return definition.parent
}

// Ref --
func (definition *FilterDefinition) Ref(ref string) *RouteDefinition {
	definition.parent.AddFactory(func(context *Context, parent Processor) (Processor, Service, error) {
		registry := context.Registry()
		ifc, err := registry.Lookup(ref)

		if ifc != nil && err == nil {
			if predicate, ok := ifc.(func(e *Exchange) bool); ok {
				fn := func(e *Exchange, out chan<- *Exchange) {
					if predicate(e) {
						out <- e
					}
				}

				p := NewProcessor(fn)
				p.Parent(parent)

				return p, nil, nil
			}
		}

		// TODO: error handling
		return parent, nil, fmt.Errorf("Unsupported type for ref:%s, type=%T", ref, ifc)
	})

	return definition.parent
}
