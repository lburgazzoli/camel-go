package camel

import "github.com/lburgazzoli/camel-go/api"

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
func (definition *ToDefinition) Unwrap(context *Context, parent Processor) (Processor, api.Service, error) {
	var err error
	var producer Producer
	var endpoint Endpoint

	if endpoint, err = context.CreateEndpointFromURI(definition.URI); err != nil {
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
