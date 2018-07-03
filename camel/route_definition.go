package camel

import "github.com/lburgazzoli/camel-go/api"

// ==========================
//
//
//
// ==========================

// Unwrappable --
type Unwrappable interface {
	Unwrap(context *Context, parent Processor) (Processor, api.Service, error)
}

// ==========================
//
//
//
// ==========================

// Definition --
type Definition interface {
	Parent() Definition
	Children() []Definition
}

// ==========================
//
//
//
// ==========================

// RouteDefinition --
type RouteDefinition struct {
	parent   Definition
	children []Definition
}

// Parent --
func (definition *RouteDefinition) Parent() Definition {
	return definition.parent
}

// Children --
func (definition *RouteDefinition) Children() []Definition {
	return definition.children
}

// AddChild --
func (definition *RouteDefinition) AddChild(child Definition) *RouteDefinition {
	if definition.children == nil {
		definition.children = make([]Definition, 0)
	}

	definition.children = append(definition.children, child)

	return definition
}

// ==========================
//
//
//
// ==========================

// From --
func From(uri string) *RouteDefinition {
	from := FromDefinition{}
	from.parent = nil
	from.children = nil
	from.URI = uri

	def := RouteDefinition{}
	def.parent = &from
	def.children = make([]Definition, 0)

	from.children = []Definition{&def}

	return &def
}
