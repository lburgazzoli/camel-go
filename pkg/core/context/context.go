package context

import (
	"context"
	"io"
	"time"

	"go.uber.org/zap"

	"github.com/lburgazzoli/camel-go/pkg/core/typeconverter"

	"github.com/lburgazzoli/camel-go/pkg/core/properties"

	"github.com/pkg/errors"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"

	"github.com/lburgazzoli/camel-go/pkg/core/route"

	"github.com/lburgazzoli/camel-go/pkg/core/registry"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

func NewDefaultContext(logger *zap.Logger) camel.Context {
	p, err := properties.NewDefaultProperties()
	if err != nil {
		// TODO: must return an error
		panic(err)
	}

	r, err := registry.NewDefaultRegistry()
	if err != nil {
		// TODO: must return an error
		panic(err)
	}

	tc, err := typeconverter.NewDefaultTypeConverter()
	if err != nil {
		// TODO: must return an error
		panic(err)
	}

	id := uuid.New()

	ctx := defaultContext{
		id:            id,
		system:        actor.NewActorSystem(),
		registry:      r,
		properties:    p,
		typeConverter: tc,
		verticles:     make(map[string]vh),
		logger:        logger.With(zap.String("context.id", id)),
	}

	return &ctx
}

type vh struct {
	V camel.Verticle
	P *actor.PID
}

type defaultContext struct {
	id            string
	system        *actor.ActorSystem
	registry      camel.Registry
	properties    camel.Properties
	typeConverter camel.TypeConverter
	verticles     map[string]vh
	logger        *zap.Logger
}

func (c *defaultContext) ID() string {
	return c.id
}

func (c *defaultContext) Start(context.Context) error {
	return nil
}

func (c *defaultContext) Stop(context.Context) error {
	return nil
}

func (c *defaultContext) Close(context.Context) error {
	return nil
}

func (c *defaultContext) LoadRoutes(ctx context.Context, in io.Reader) error {
	routes, err := route.Load(in)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, camel.ContextKeyCamelContext, c)

	for i := range routes {
		if _, err := routes[i].Reify(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (c *defaultContext) Registry() camel.Registry {
	return c.registry
}

func (c *defaultContext) Spawn(v camel.Verticle) error {
	f := func() actor.Actor { return v }
	p := actor.PropsFromProducer(f)

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

func (c *defaultContext) Send(id string, message camel.Message) error {
	v, ok := c.verticles[id]
	if !ok {
		return camelerrors.NotFoundf("verticle with id %s not found in registry", id)
	}

	c.system.Root.Send(v.P, message)

	return nil
}

func (c *defaultContext) Receive(_ string, _ time.Duration) (camel.Message, error) {
	return nil, camelerrors.NotImplemented("Receive")
}

func (c *defaultContext) Properties() camel.Properties {
	return c.properties
}

func (c *defaultContext) TypeConverter() camel.TypeConverter {
	return c.typeConverter
}

func (c *defaultContext) Logger() *zap.Logger {
	return c.logger
}
