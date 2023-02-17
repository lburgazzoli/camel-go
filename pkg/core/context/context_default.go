package context

import (
	"io"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

func NewDefaultContext() Context {
	ctx := defaultContext{
		id:       uuid.New(),
		system:   actor.NewActorSystem(),
		registry: NewDefaultRegistry(),
	}

	return &ctx
}

type defaultContext struct {
	id       string
	system   *actor.ActorSystem
	registry Registry
}

func (c *defaultContext) ID() string {
	return c.id
}

func (c *defaultContext) LoadRoutes(in io.Reader) error {
	return nil
}

func (c *defaultContext) Registry() Registry {
	return c.registry
}
