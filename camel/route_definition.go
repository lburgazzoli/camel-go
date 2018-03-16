package camel

import (
	"errors"

	"github.com/rs/zerolog/log"
)

// ==========================
//
// ==========================

type definitionFactory func(parent *Pipe) (*Pipe, Service)

// ==========================
//
// RouteDefinition
//
//    WORK IN PROGRESS
//
// ==========================

// NewRouteDefinition --
func NewRouteDefinition(context *Context) *RouteDefinition {
	return &RouteDefinition{context: context}
}

// RouteDefinition --
type RouteDefinition struct {
	context             *Context
	definition          definitionFactory
	processorDefinition *ProcessorDefinition
}

func (definition *RouteDefinition) addDefinitionsToRoute(route *Route, rootPipe *Pipe, rootDefinition *ProcessorDefinition) {
	var s Service

	if rootDefinition.definitions != nil {
		for _, def := range rootDefinition.definitions {
			rootPipe, s = def(rootPipe)

			route.AddService(s)
		}
	}

	if rootDefinition.child != nil {
		definition.addDefinitionsToRoute(route, rootPipe, rootDefinition.child)
	}
}

// ToRoute --
func (definition *RouteDefinition) ToRoute(context *Context) (*Route, error) {
	route := Route{}

	if definition.definition != nil {
		p, s := definition.definition(nil)
		route.AddService(s)

		if definition.processorDefinition != nil {
			definition.addDefinitionsToRoute(&route, p, definition.processorDefinition)
		}
	} else {
		return nil, errors.New("no from")
	}

	return &route, nil
}

// From --
func (definition *RouteDefinition) From(uri string) *ProcessorDefinition {
	var err error
	var consumer Consumer
	var endpoint Endpoint

	if endpoint, err = definition.context.CreateEndpointFromURI(uri); err != nil {
		return nil
	}

	if consumer, err = endpoint.CreateConsumer(); err != nil {
		return nil
	}

	definition.definition = func(parent *Pipe) (*Pipe, Service) {
		if parent != nil {
			log.Panic().Msgf("parent pipe should be nil, got %+v", parent)
		}

		return consumer.Pipe(), consumer
	}

	definition.processorDefinition = &ProcessorDefinition{}
	definition.processorDefinition.parent = nil
	definition.processorDefinition.child = nil
	definition.processorDefinition.context = definition.context
	definition.processorDefinition.definitions = make([]definitionFactory, 0)

	return definition.processorDefinition
}
