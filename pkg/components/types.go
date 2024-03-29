package components

import (
	"log/slog"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
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
		logger: ctx.Logger().
			With(slog.String("component.scheme", scheme), slog.String("component.id", id)),
	}

	return dc
}

type DefaultComponent struct {
	ctx    api.Context
	scheme string
	id     string
	logger *slog.Logger
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

func (c *DefaultComponent) Logger() *slog.Logger {
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
		logger: component.Logger().
			With(slog.String("endpoint.id", id)),
	}
}

type DefaultEndpoint struct {
	component api.Component
	id        string
	logger    *slog.Logger
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

func (e *DefaultEndpoint) Logger() *slog.Logger {
	return e.logger
}

//
// Default Consumer
//

func NewDefaultConsumer(endpoint api.Endpoint, target *actor.PID) DefaultConsumer {
	v := processors.NewDefaultVerticle()

	return DefaultConsumer{
		DefaultVerticle: v,
		target:          target,
		endpoint:        endpoint,
		logger: endpoint.Logger().
			With(slog.String("consumer.id", v.ID())),
	}
}

type DefaultConsumer struct {
	processors.DefaultVerticle

	logger   *slog.Logger
	endpoint api.Endpoint
	target   *actor.PID
}

func (c *DefaultConsumer) Logger() *slog.Logger {
	return c.logger
}

func (c *DefaultConsumer) Target() *actor.PID {
	return c.target
}

func (c *DefaultConsumer) Endpoint() api.Endpoint {
	return c.endpoint
}

//
// Default Producer
//

func NewDefaultProducer(endpoint api.Endpoint) DefaultProducer {
	v := processors.NewDefaultVerticle()

	return DefaultProducer{
		DefaultVerticle: v,
		endpoint:        endpoint,
		logger: endpoint.Logger().
			With(slog.String("producer.id", v.ID())),
	}
}

type DefaultProducer struct {
	processors.DefaultVerticle

	logger   *slog.Logger
	endpoint api.Endpoint
}

func (p *DefaultProducer) Logger() *slog.Logger {
	return p.logger
}

func (p *DefaultProducer) Endpoint() api.Endpoint {
	return p.endpoint
}
