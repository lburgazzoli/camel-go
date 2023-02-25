package components

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

type ComponentFactory func(api.Context, map[string]interface{}) (api.Component, error)

var Factories = make(map[string]ComponentFactory)

//
// Default Component
//

func NewDefaultComponent(ctx api.Context, scheme string) DefaultComponent {
	return DefaultComponent{
		ctx:    ctx,
		scheme: scheme,
		id:     uuid.New(),
	}
}

type DefaultComponent struct {
	ctx    api.Context
	scheme string
	id     string
}

func (c *DefaultComponent) Context() api.Context {
	return c.ctx
}

func (c *DefaultComponent) ID() string {
	return c.id
}

func (c *DefaultComponent) Scheme() string {
	return c.scheme
}

//
// Default Endpoint
//

func NewDefaultEndpoint(component api.Component) DefaultEndpoint {
	return DefaultEndpoint{
		component: component,
		id:        uuid.New(),
	}
}

type DefaultEndpoint struct {
	component api.Component
	id        string
}

func (e *DefaultEndpoint) Component() api.Component {
	return e.component
}

func (e *DefaultEndpoint) ID() string {
	return e.id
}
