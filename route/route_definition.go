package route

import (
	"github.com/lburgazzoli/camel-go/api"
	zlog "github.com/rs/zerolog/log"
)

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

// ToRoute --
func ToRoute(context api.Context, definition Definition) (*api.Route, error) {
	route := api.NewRoute("")

	// Find the root
	for definition.Parent() != nil {
		definition = definition.Parent()
	}

	unwrapDefinition(context, route, nil, definition)

	return nil, nil
}

func unwrapDefinition(context api.Context, route *api.Route, processor api.Processor, definition Definition) api.Processor {
	var s api.Service
	var e error

	p := processor

	if u, ok := definition.(api.Unwrappable); ok {
		p, s, e = u.Unwrap(context, p)

		if e != nil {
			zlog.Fatal().Msgf("unable to load processor %v (%s)", definition, e)
		}

		if e == nil && s != nil {
			route.AddService(s)
		}

		if p != nil {
			if s, ok := p.(api.Service); ok {
				route.AddService(s)
			}
		} else {
			p = processor
		}
	}

	for _, c := range definition.Children() {
		p = unwrapDefinition(context, route, p, c)
	}

	return p
}
