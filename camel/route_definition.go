package camel

import (
	"errors"

	"github.com/rs/zerolog/log"
)

// ProcessorDefinition --
type ProcessorDefinition interface {
	To(uri string) ProcessorDefinition
	Process(processor Processor) ProcessorDefinition
}

// RouteDefinition --
type RouteDefinition interface {
	ToRoute

	From(uri string) ProcessorDefinition
}

// NewRouteDefinition --
func NewRouteDefinition(context *Context) RouteDefinition {
	return &defaultRouteDefinition{context: context}
}

// ==========================
//
// ==========================

type definitionFactory func(parent *Pipe) (*Pipe, Service)

// ==========================
//
// RouteDefinition Impl
//
//    WORK IN PROGRESS
//
// ==========================

type defaultRouteDefinition struct {
	RouteDefinition

	context             *Context
	definition          definitionFactory
	processorDefinition *defaultProcessorDefinition
}

// ToRoute --
func (definition *defaultRouteDefinition) ToRoute(context *Context) (*Route, error) {
	route := Route{}

	if definition.definition != nil {
		var p *Pipe
		var s Service

		p, s = definition.definition(nil)
		route.AddService(s)

		if definition.processorDefinition != nil {
			for _, def := range definition.processorDefinition.definitions {
				p, s = def(p)

				route.AddService(s)
			}
		}
	} else {
		return nil, errors.New("No from")
	}

	return &route, nil
}

// From --
func (definition *defaultRouteDefinition) From(uri string) ProcessorDefinition {
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
			log.Panic().Msgf("Parent pipe should be nil, got %+v", parent)
		}

		return consumer.Pipe(), consumer
	}

	definition.processorDefinition = &defaultProcessorDefinition{}
	definition.processorDefinition.parent = nil
	definition.processorDefinition.child = nil
	definition.processorDefinition.context = definition.context
	definition.processorDefinition.definitions = make([]definitionFactory, 0)

	return definition.processorDefinition
}

// ==========================
//
// ProcessorDefinition Impl
//
//    WORK IN PROGRESS
//
// ==========================

type defaultProcessorDefinition struct {
	ProcessorDefinition

	context     *Context
	definitions []definitionFactory
	child       *defaultProcessorDefinition
	parent      *defaultProcessorDefinition
}

func (definition *defaultProcessorDefinition) addFactory(factory definitionFactory) ProcessorDefinition {
	definition.definitions = append(definition.definitions, factory)

	return definition
}

// To --
func (definition *defaultProcessorDefinition) To(uri string) ProcessorDefinition {
	var err error
	var producer Producer
	var endpoint Endpoint

	if endpoint, err = definition.context.CreateEndpointFromURI(uri); err != nil {
		return nil
	}

	if producer, err = endpoint.CreateProducer(); err != nil {
		return nil
	}

	return definition.addFactory(func(parent *Pipe) (*Pipe, Service) {
		p := producer.Pipe()
		parent.Next(p)

		return p, producer
	})
}

func (definition *defaultProcessorDefinition) Process(processor Processor) ProcessorDefinition {
	return definition.addFactory(func(parent *Pipe) (*Pipe, Service) {
		return parent.Process(processor), nil
	})
}
