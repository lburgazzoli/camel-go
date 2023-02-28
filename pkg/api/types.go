package api

import (
	"context"
	"io"
	"time"

	"github.com/asynkron/protoactor-go/actor"

	ce "github.com/cloudevents/sdk-go/v2"
)

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
}

type Context interface {
	Identifiable
	Service
	Closer

	Registry() Registry
	LoadRoutes(ctx context.Context, in io.Reader) error

	// Spawn ---
	// TODO: must be hidden
	// TODO: each route must have its own context/supervisor
	Spawn(Verticle) error

	// Send ---
	// TODO: must use name instead of PID
	Send(string, Message) error

	// Receive ---
	// TODO: must use name instead of PID
	Receive(string, time.Duration) (Message, error)
}

type Component interface {
	Identifiable

	Context() Context
	Scheme() string
	Endpoint(Parameters) (Endpoint, error)
}

type Endpoint interface {
	Identifiable
	Service

	Component() Component
}

type Message interface {
	ce.EventContext

	Fail(error)
	Error() error

	// Annotation ---
	// TODO: add options Annotation("foo", opt.WithDefault("bar"), opt.AsType(baz{})).
	Annotation(string) (string, bool)
	SetAnnotation(string, string)

	// Content ---
	// TODO: add options Content(opt.AsType(baz{})).
	Content() interface{}
	SetContent(interface{})
}

type Processor = func(Message)

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

// OutputAware ---
// TODO: use name or other abstractions instead of PIO.
type OutputAware interface {
	Next(string)
	Outputs() []string
}

// WithOutputs ---
// TODO: move to helper package.
type WithOutputs struct {
	outputs []string
}

func (o *WithOutputs) Next(id string) {
	if id == "" {
		return
	}

	o.outputs = append(o.outputs, id)
}

func (o *WithOutputs) Outputs() []string {
	return o.outputs
}

type Verticle interface {
	Identifiable
	OutputAware

	actor.Actor
}
