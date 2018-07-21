// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

			if s, ok := p.(api.Service); ok {
				route.AddService(s)
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
