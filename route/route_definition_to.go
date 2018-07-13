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
	parent   *RouteDefinition
	children []Definition

	URI string
}

// Parent --
func (definition *ToDefinition) Parent() Definition {
	return definition.parent
}

// Children --
func (definition *ToDefinition) Children() []Definition {
	return definition.children
}

// Unwrap ---
func (definition *ToDefinition) Unwrap(context api.Context, parent api.Processor) (api.Processor, api.Service, error) {
	var err error
	var producer api.Producer
	var endpoint api.Endpoint

	if endpoint, err = camel.NewEndpointFromURI(context, definition.URI); err != nil {
		return parent, nil, err
	}

	if producer, err = endpoint.CreateProducer(); err != nil {
		return parent, nil, err
	}
	p := producer.Processor()

	parent.Subscribe(func(e api.Exchange) {
		p.Publish(e)
	})

	return producer.Processor(), producer, nil
}
