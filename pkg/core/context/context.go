package context

import (
	"io"
	"time"

	"github.com/pkg/errors"

	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"

	"github.com/lburgazzoli/camel-go/pkg/core/route"

	"github.com/lburgazzoli/camel-go/pkg/core/registry"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

func NewDefaultContext() api.Context {
	ctx := defaultContext{
		id:        uuid.New(),
		system:    actor.NewActorSystem(),
		registry:  registry.NewDefaultRegistry(),
		verticles: make(map[string]vh),
	}

	return &ctx
}

type vh struct {
	V api.Verticle
	P *actor.PID
}

type defaultContext struct {
	id        string
	system    *actor.ActorSystem
	registry  api.Registry
	verticles map[string]vh
}

func (c *defaultContext) ID() string {
	return c.id
}

func (c *defaultContext) LoadRoutes(in io.Reader) error {
	routes, err := route.Load(in)
	if err != nil {
		return err
	}

	for i := range routes {
		if _, err := routes[i].Reify(c); err != nil {
			return err
		}
	}

	return nil
}

func (c *defaultContext) Registry() api.Registry {
	return c.registry
}

func (c *defaultContext) Spawn(v api.Verticle) error {
	p := actor.PropsFromProducer(func() actor.Actor {
		return v
	})

	pid, err := c.system.Root.SpawnNamed(p, v.ID())
	if err != nil {
		return errors.Wrapf(err, "unable to spawn verticle with id %s", v.ID())
	}

	c.verticles[v.ID()] = vh{
		V: v,
		P: pid,
	}

	return nil
}

func (c *defaultContext) Send(id string, message api.Message) error {
	v, ok := c.verticles[id]
	if !ok {
		return camelerrors.NotFoundf("verticle with id %s not found in registry", id)
	}

	c.system.Root.Send(v.P, message)

	return nil
}

func (c *defaultContext) Receive(_ string, _ time.Duration) (api.Message, error) {
	return nil, camelerrors.NotImplemented("Receive")
}
