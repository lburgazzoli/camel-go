package camel

import (
	"fmt"
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
	parent   *RouteDefinition
	children []Definition

	predicate    func(*Exchange) bool
	predicateRef string
}

// Parent --
func (definition *FilterDefinition) Parent() Definition {
	return definition.parent
}

// Children --
func (definition *FilterDefinition) Children() []Definition {
	return definition.children
}

// Unwrap ---
func (definition *FilterDefinition) Unwrap(context *Context, parent Processor) (Processor, Service, error) {
	if definition.predicate != nil {
		p := NewProcessorWithParent(parent, func(e *Exchange, out chan<- *Exchange) {
			if definition.predicate(e) {
				out <- e
			}
		})

		return p, nil, nil
	}

	if definition.predicateRef != "" {
		registry := context.Registry()
		ifc, err := registry.Lookup(definition.predicateRef)

		if ifc != nil && err == nil {
			if predicate, ok := ifc.(func(e *Exchange) bool); ok {
				p := NewProcessorWithParent(parent, func(e *Exchange, out chan<- *Exchange) {
					if predicate(e) {
						out <- e
					}
				})

				return p, nil, nil
			}
		}

		// TODO: error handling
		return nil, nil, fmt.Errorf("Unsupported type for ref:%s, type=%T", definition.predicateRef, ifc)
	}

	return nil, nil, nil

}

// Fn --
func (definition *FilterDefinition) Fn(predicate func(*Exchange) bool) *RouteDefinition {
	definition.predicate = predicate
	return definition.parent
}

// Ref --
func (definition *FilterDefinition) Ref(ref string) *RouteDefinition {
	definition.predicateRef = ref
	return definition.parent
}
