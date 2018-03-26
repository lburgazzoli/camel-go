package camel

import "github.com/rs/zerolog/log"

// ==========================
//
// ==========================

// DefinitionFactory --
type DefinitionFactory func(context *Context, parent *Subject) (*Subject, Service, error)

// ==========================
//
// RouteDefinition
//
//    WORK IN PROGRESS
//
// ==========================

// From --
func From(uri string) *RouteDefinition {
	definition := RouteDefinition{factories: make([]DefinitionFactory, 0)}
	definition.AddFactory(func(context *Context, parent *Subject) (*Subject, Service, error) {
		var err error
		var consumer Consumer
		var endpoint Endpoint

		if endpoint, err = context.CreateEndpointFromURI(uri); err != nil {
			return parent, nil, nil
		}

		if consumer, err = endpoint.CreateConsumer(); err != nil {
			return parent, nil, nil
		}

		if parent != nil {
			log.Panic().Msgf("parent pipe should be nil, got %+v", parent)
		}

		return consumer.Subject(), consumer, nil
	})

	return &definition
}

// RouteDefinition --
type RouteDefinition struct {
	factories []DefinitionFactory
	child     *RouteDefinition
	parent    *RouteDefinition
}

// AddFactory --
func (definition *RouteDefinition) AddFactory(factory DefinitionFactory) *RouteDefinition {
	definition.factories = append(definition.factories, factory)

	return definition
}

// End --
func (definition *RouteDefinition) End() *RouteDefinition {
	return definition.parent
}

// To --
func (definition *RouteDefinition) To(uri string) *RouteDefinition {
	return definition.AddFactory(func(context *Context, parent *Subject) (*Subject, Service, error) {
		var err error
		var producer Producer
		var endpoint Endpoint

		if endpoint, err = context.CreateEndpointFromURI(uri); err != nil {
			return parent, nil, err
		}

		if producer, err = endpoint.CreateProducer(); err != nil {
			return parent, nil, err
		}
		p := producer.Subject()

		parent.Subscribe(func(e *Exchange) {
			p.Publish(e)
		})

		return p, producer, nil
	})
}
