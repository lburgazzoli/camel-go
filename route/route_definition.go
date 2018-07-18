package route

import (
	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/processor"
	zlog "github.com/rs/zerolog/log"
)

// ==========================
//
//
//
// ==========================

// ProcessingNode --
type ProcessingNode interface {
	Processor() (api.Processor, error)
}

// ServiceNode --
type ServiceNode interface {
	// TODO: refactor
	Service() (api.Processor, api.Service, error)
}

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

	if p := unwrapDefinition(context, route, nil, definition); p != nil {
		p.Subscribe(func(_ api.Exchange) {
			// processing end
		})
	}

	return route, nil
}

func unwrapDefinition(context api.Context, route *api.Route, parent api.Processor, definition Definition) api.Processor {
	var p api.Processor
	var s api.Service
	var e error

	p = parent

	if node, ok := definition.(api.ContextAware); ok {
		node.SetContext(context)
	}

	if node, ok := definition.(ProcessingNode); ok {
		p, e = node.Processor()

		if e != nil {
			zlog.Fatal().Msgf("unable to load processing node %v (%s)", definition, e)
		}

		if p != nil {
			if parent != nil {
				zlog.Debug().Msgf("connect %+v", definition)
				processor.Connect(parent, p)
			}
		} else {
			p = parent
		}
	}

	if node, ok := definition.(ServiceNode); ok {
		p, s, e = node.Service()

		if e != nil {
			zlog.Fatal().Msgf("unable to load service node %v (%s)", definition, e)
		}

		if s != nil {
			route.AddService(s)
		}

		if p != nil {
			if parent != nil {
				zlog.Debug().Msgf("connect %+v, %+v", parent, definition)
				processor.Connect(parent, p)
			}
		} else {
			p = parent
		}
	}

	for _, c := range definition.Children() {
		p = unwrapDefinition(context, route, p, c)
	}

	return p
}
