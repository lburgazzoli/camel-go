package components

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"go.uber.org/zap"
)

type ComponentFactory func(api.Context, map[string]interface{}) (api.Component, error)

var Factories = make(map[string]ComponentFactory)

//
// Default Component
//

func NewDefaultComponent(ctx api.Context, scheme string) DefaultComponent {
	id := uuid.New()

	dc := DefaultComponent{
		ctx:    ctx,
		scheme: scheme,
		id:     id,
		logger: ctx.Logger().With(zap.String("component.scheme", scheme), zap.String("component.id", id)),
	}

	return dc
}

type DefaultComponent struct {
	ctx    api.Context
	scheme string
	id     string
	logger *zap.Logger
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

func (c *DefaultComponent) Logger() *zap.Logger {
	return c.logger
}

//
// Default Endpoint
//

func NewDefaultEndpoint(component api.Component) DefaultEndpoint {
	id := uuid.New()

	return DefaultEndpoint{
		component: component,
		id:        id,
		logger:    component.Logger().With(zap.String("endpoint.id", id)),
	}
}

type DefaultEndpoint struct {
	component api.Component
	id        string
	logger    *zap.Logger
}

func (e *DefaultEndpoint) Context() api.Context {
	return e.component.Context()
}

func (e *DefaultEndpoint) Component() api.Component {
	return e.component
}

func (e *DefaultEndpoint) ID() string {
	return e.id
}

func (e *DefaultEndpoint) Logger() *zap.Logger {
	return e.logger
}

//
// Default Consumer
//

func NewDefaultConsumer(endpoint api.Endpoint, target *actor.PID) DefaultConsumer {
	id := uuid.New()

	return DefaultConsumer{
		target: target,
		logger: endpoint.Logger().With(zap.String("consumer.id", id)).Sugar(),
	}
}

type DefaultConsumer struct {
	processors.DefaultVerticle

	logger *zap.SugaredLogger
	target *actor.PID
}

func (c *DefaultConsumer) Logger() *zap.SugaredLogger {
	return c.logger
}

func (c *DefaultConsumer) Target() *actor.PID {
	return c.target
}
