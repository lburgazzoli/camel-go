package api

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.uber.org/zap"

	"github.com/asynkron/protoactor-go/actor"

	ce "github.com/cloudevents/sdk-go/v2"
)

type ContextKey string

const (
	ContextKeyCamelContext = ContextKey("camel-context")
	ContextKeyActorContext = ContextKey("actor-context")
)

func GetContext(ctx context.Context) Context {
	value := ctx.Value(ContextKeyCamelContext)
	if value == nil {
		panic(fmt.Errorf("unable to get CamelContext from context"))
	}

	answer, ok := value.(Context)
	if !ok {
		panic(fmt.Errorf("type cast error %v", value))
	}

	return answer
}

func GetActorContext(ctx context.Context) actor.Context {
	value := ctx.Value(ContextKeyActorContext)
	if value == nil {
		panic(fmt.Errorf("unable to get actor Context from context"))
	}

	answer, ok := value.(actor.Context)
	if !ok {
		panic(fmt.Errorf("type cast error %v", value))
	}

	return answer
}

type Parameters map[string]interface{}

type Closer interface {
	// Close closes the resource.
	Close(context.Context) error
}

type Service interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type Identifiable interface {
	ID() string
}

type Registry interface {
	Get(string) (interface{}, bool)
	Set(string, interface{})
	Del(key string) interface{}
}

type Properties interface {
	AddSource(string) error
	String(string) string
}

//nolint:interfacebloat
type Context interface {
	Identifiable
	Service
	Closer

	Registry() Registry
	Properties() Properties
	TypeConverter() TypeConverter

	LoadRoutes(ctx context.Context, in io.Reader) error

	// Spawn ---
	// TODO: must be hidden
	// TODO: each route must have its own context/supervisor
	Spawn(Verticle) (*actor.PID, error)

	// Send ---
	Send(string, Message) error

	// SendTo ---
	// TODO: must be hidden maybe
	SendTo(*actor.PID, Message) error

	// Receive ---
	Receive(string, time.Duration) (Message, error)

	Logger() *zap.Logger
}

type Component interface {
	Identifiable

	Context() Context
	Scheme() string
	Endpoint(Parameters) (Endpoint, error)

	Logger() *zap.Logger
}

type Endpoint interface {
	Identifiable
	Service

	Context() Context
	Component() Component

	Logger() *zap.Logger
}

//nolint:interfacebloat
type Message interface {
	ce.EventContext

	Fail(error)

	SetError(error)
	Error() error

	// Annotation ---
	// TODO: add options Annotation("foo", opt.WithDefault("bar"), opt.AsType(baz{})).
	Annotation(string) (string, bool)
	SetAnnotation(string, string)

	Annotations() map[string]string
	SetAnnotations(map[string]string)
	ForEachAnnotation(func(string, string))

	// Content ---
	// TODO: add options Content(opt.AsType(baz{})).
	Content() interface{}
	SetContent(interface{})
}

type Processor = func(context.Context, Message) error
type Predicate = func(context.Context, Message) (bool, error)

type Producer interface {
	Service
	Verticle

	Endpoint() Endpoint
}

type ProducerFactory interface {
	Producer() (Producer, error)
}

type Consumer interface {
	Service
	Verticle

	Endpoint() Endpoint
}

type ConsumerFactory interface {
	Consumer() (Consumer, error)
}

type OutputAware interface {
	Next(*actor.PID)
	Outputs() []*actor.PID
}

type WithOutputs struct {
	outputs []*actor.PID
}

func (o *WithOutputs) Next(id *actor.PID) {
	if id == nil {
		return
	}

	o.outputs = append(o.outputs, id)
}

func (o *WithOutputs) Outputs() []*actor.PID {
	return o.outputs
}

type Verticle interface {
	Identifiable
	OutputAware

	actor.Actor
}

type TypeConverter interface {
	Convert(interface{}, interface{}) (bool, error)
}

type RawJSON map[string]interface{}
