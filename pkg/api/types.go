package api

import (
	"context"
	"io"

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
