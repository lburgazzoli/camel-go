package api

import (
	"context"
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

	C() context.Context

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

	Scheme() string
	Endpoint(Parameters) (Endpoint, error)
}

type Endpoint interface {
	Identifiable
	Service
}

type Message interface {
	ce.EventContext

	Fail(error)
	Error() error

	Annotation(string) (interface{}, bool)
	SetAnnotation(string, interface{})

	Content() interface{}
	SetContent(interface{})
}

type Processor func(Message)
