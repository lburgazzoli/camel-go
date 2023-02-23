package context

import (
	"context"
	"io"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

func NewDefaultContext(context context.Context) api.Context {
	ctx := defaultContext{
		ctx:      context,
		id:       uuid.New(),
		system:   actor.NewActorSystem(),
		registry: NewDefaultRegistry(),
	}

	return &ctx
}

type defaultContext struct {
	ctx      context.Context
	id       string
	system   *actor.ActorSystem
	registry api.Registry
}

func (c *defaultContext) C() context.Context {
	return c.ctx
}

func (c *defaultContext) ID() string {
	return c.id
}

func (c *defaultContext) LoadRoutes(in io.Reader) error {
	return nil
}

func (c *defaultContext) Registry() api.Registry {
	return c.registry
}
