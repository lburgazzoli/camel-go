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
		return api.NewFilteringPipeline(definition.predicate), nil
	}

	if definition.predicateRef != "" {
		registry := definition.context.Registry()
		ifc, found := registry.Lookup(definition.predicateRef)

		if ifc != nil && found {
			if p, ok := ifc.(func(e api.Exchange) bool); ok {
				return api.NewFilteringPipeline(p), nil
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
