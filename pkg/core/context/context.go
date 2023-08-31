package context

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"time"

	"github.com/lburgazzoli/camel-go/pkg/core/message"

	"github.com/lburgazzoli/camel-go/pkg/core/errors/log"

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

func NewDefaultContext(logger *slog.Logger, opts ...Option) camel.Context {
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
		logger:        logger.With(slog.String("context.id", id)),
	}

	for _, opt := range opts {
		opt(&ctx)
	}

	if ctx.errorHandler != nil {
		ctx.system.Root.WithGuardian(ctx.errorHandler)
	}

	return &ctx
}

type vh struct {
	V camel.Verticle
	P *actor.PID
}

type Option func(c *defaultContext)

func WithErrorHandler(handler camelerrors.Handler) Option {
	return func(c *defaultContext) {
		c.errorHandler = handler
	}
}
func WithLogErrorHandler() Option {
	return func(c *defaultContext) {
		c.errorHandler = log.NewLogHandler(c)
	}
}

type defaultContext struct {
	id            string
	system        *actor.ActorSystem
	registry      camel.Registry
	properties    camel.Properties
	typeConverter camel.TypeConverter
	verticles     map[string]vh
	logger        *slog.Logger
	errorHandler  camelerrors.Handler
}

func (c *defaultContext) NewMessage() camel.Message {
	return message.New(c)
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
		v, err := routes[i].Reify(ctx)
		if err != nil {
			return err
		}

		_, err = c.Spawn(v)
		if err != nil {
			panic(errors.Wrapf(err, "unable to spawn route verticle with id %s", v.ID()))
		}
	}

	return nil
}

func (c *defaultContext) Registry() camel.Registry {
	return c.registry
}

func (c *defaultContext) Spawn(v camel.Verticle) (*actor.PID, error) {
	f := func() actor.Actor { return v }
	p := actor.PropsFromProducer(f)

	pid, err := c.system.Root.SpawnNamed(p, v.ID())
	if err != nil {
		return nil, errors.Wrapf(err, "unable to spawn verticle with id %s", v.ID())
	}

	c.verticles[v.ID()] = vh{
		V: v,
		P: pid,
	}

	return pid, nil
}

func (c *defaultContext) Send(id string, message camel.Message) error {
	target, ok := registry.GetAs[*actor.PID](c.Registry(), id)
	if !ok {
		return camelerrors.NotFoundf("verticle with id %s not found in registry", id)
	}

	return c.SendTo(target, message)
}

func (c *defaultContext) SendTo(target *actor.PID, message camel.Message) error {
	c.system.Root.Request(target, message)
	return nil
}

func (c *defaultContext) SendToAs(target *actor.PID, sender *actor.PID, message camel.Message) error {
	c.system.Root.RequestWithCustomSender(target, message, sender)
	return nil
}

func (c *defaultContext) Receive(_ string, _ time.Duration) (camel.Message, error) {
	return nil, camelerrors.NotImplemented("Receive")
}

func (c *defaultContext) RequestTo(target *actor.PID, message camel.Message, timeout time.Duration) (camel.Message, error) {
	r, err := c.system.Root.RequestFuture(target, message, timeout).Result()
	if err != nil {
		return nil, err
	}

	m, ok := r.(camel.Message)
	if !ok {
		return nil, fmt.Errorf("response is not of type camel.Message, got %v", reflect.TypeOf(r))
	}

	return m, nil
}

func (c *defaultContext) Properties() camel.Properties {
	return c.properties
}

func (c *defaultContext) TypeConverter() camel.TypeConverter {
	return c.typeConverter
}

func (c *defaultContext) Logger() *slog.Logger {
	return c.logger
}
