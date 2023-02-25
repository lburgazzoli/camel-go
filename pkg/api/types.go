package api

import (
	"io"
	"time"

	"github.com/asynkron/protoactor-go/actor"

	ce "github.com/cloudevents/sdk-go/v2"
)

type Parameters map[string]interface{}

type Service interface {
	Start() error
	Stop() error
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

	Registry() Registry
	LoadRoutes(in io.Reader) error

	// Spawn ---
	// TODO: must be hidden
	// TODO: each route must have its own context/supervisor
	Spawn(string, actor.Actor) (*actor.PID, error)

	// SpawnFn ---
	// TODO: must be hidden
	// TODO: each route must have its own context/supervisor
	SpawnFn(string, actor.ReceiveFunc) (*actor.PID, error)

	// Send ---
	// TODO: must use name instead of PID
	Send(*actor.PID, Message)

	// Receive ---
	// TODO: must use name instead of PID
	Receive(*actor.PID, time.Duration) Message
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
	Annotation(string) (interface{}, bool)
	SetAnnotation(string, interface{})

	// Content ---
	// TODO: add options Content(opt.AsType(baz{})).
	Content() interface{}
	SetContent(interface{})
}

type Processor = func(Message)

type Producer interface {
	Service
	OutputAware
	
	actor.Actor

	Endpoint() Endpoint
}

type ProducerFactory interface {
	Producer() (Producer, error)
}

type Consumer interface {
	Service
	OutputAware

	actor.Actor

	Endpoint() Endpoint
}

type ConsumerFactory interface {
	Consumer() (Consumer, error)
}

// OutputAware ---
// TODO: use name or other abstractions instead of PIO.
type OutputAware interface {
	Next(*actor.PID)
	Outputs() []*actor.PID
}

// WithOutputs ---
// TODO: move to helper package.
type WithOutputs struct {
	outputs *actor.PIDSet
}

func (o *WithOutputs) Next(pid *actor.PID) {
	if pid == nil {
		return
	}
	if o.outputs == nil {
		o.outputs = actor.NewPIDSet()
	}

	o.outputs.Add(pid)
}

func (o *WithOutputs) Outputs() []*actor.PID {
	if o.outputs == nil {
		return nil
	}

	return o.outputs.Values()
}
