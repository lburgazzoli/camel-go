package route

import (
	"fmt"

	"github.com/lburgazzoli/camel-go/api"
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

	predicate    func(api.Exchange) bool
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
func (definition *FilterDefinition) Unwrap(context api.Context, parent api.Processor) (api.Processor, api.Service, error) {
	if definition.predicate != nil {
		p := api.NewProcessorWithParent(parent, func(e api.Exchange, out chan<- api.Exchange) {
			if definition.predicate(e) {
				out <- e
			}
		})

		return p, nil, nil
	}

	if definition.predicateRef != "" {
		registry := context.Registry()
		ifc, found := registry.Lookup(definition.predicateRef)

		if ifc != nil && found {
			if predicate, ok := ifc.(func(e api.Exchange) bool); ok {
				p := api.NewProcessorWithParent(parent, func(e api.Exchange, out chan<- api.Exchange) {
					if predicate(e) {
						out <- e
					}
				})

				return p, nil, nil
			}
		}

		var err error

		if !found {
			err = fmt.Errorf("Unsupported type for ref:%s, type=%T", definition.predicateRef, ifc)
		}

		// TODO: error handling
		return nil, nil, err
	}

	return nil, nil, nil

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
