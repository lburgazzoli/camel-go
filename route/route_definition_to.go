package route

import (
	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/camel"
)

// ==========================
//
// Extend RouteDefinition DSL
//
// ==========================

// To --
func (definition *RouteDefinition) To(uri string) *RouteDefinition {
	d := ToDefinition{
		parent:   definition,
		children: nil,
		URI:      uri,
	}

	definition.AddChild(&d)

	return definition
}

// ==========================
//
// ToDefinition
//
// ==========================

// ToDefinition --
type ToDefinition struct {
	api.ContextAware
	ServiceNode

	context  api.Context
	parent   *RouteDefinition
	children []Definition

	URI string
}

// SetContext --
func (definition *ToDefinition) SetContext(context api.Context) {
	definition.context = context
}

// Context --
func (definition *ToDefinition) Context() api.Context {
	return definition.context
}

// Parent --
func (definition *ToDefinition) Parent() Definition {
	return definition.parent
}

// Children --
func (definition *ToDefinition) Children() []Definition {
	return definition.children
}

// Service ---
func (definition *ToDefinition) Service() (api.Processor, api.Service, error) {
	var err error
	var producer api.Producer
	var endpoint api.Endpoint

	if endpoint, err = camel.NewEndpointFromURI(definition.context, definition.URI); err != nil {
		return nil, nil, err
	}

	if producer, err = endpoint.CreateProducer(); err != nil {
		return nil, nil, err
	}

	return producer.Processor(), producer, nil
}
