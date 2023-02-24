package context

import (
	"context"
	"io"
	"time"

	"github.com/lburgazzoli/camel-go/pkg/core/registry"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

func NewDefaultContext(context context.Context) api.Context {
	ctx := defaultContext{
		ctx:      context,
		id:       uuid.New(),
		system:   actor.NewActorSystem(),
		registry: registry.NewDefaultRegistry(),
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

func (c *defaultContext) Spawn(name string, a actor.Actor) (*actor.PID, error) {
	p := actor.PropsFromProducer(func() actor.Actor {
		return a
	})

	return c.system.Root.SpawnNamed(p, name)
}

func (c *defaultContext) SpawnFn(name string, a actor.ReceiveFunc) (*actor.PID, error) {
	p := actor.PropsFromFunc(a)

	return c.system.Root.SpawnNamed(p, name)
}

func (c *defaultContext) Send(pid *actor.PID, message api.Message) {
	c.system.Root.Send(pid, message)
}

func (c *defaultContext) Receive(*actor.PID, time.Duration) api.Message {
	return nil
}
