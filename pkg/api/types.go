package api

import "io"

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
}

type Component interface {
	Identifiable
	
	Scheme() string
}
