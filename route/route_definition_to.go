package route

import (
	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/camel"
	"github.com/lburgazzoli/camel-go/processor"
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
	ProcessingNode

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

// Processor ---
func (definition *ToDefinition) Processor() (api.Processor, error) {
	var err error
	var producer api.Producer
	var endpoint api.Endpoint

	if endpoint, err = camel.NewEndpointFromURI(definition.context, definition.URI); err != nil {
		return nil, err
	}

	if producer, err = endpoint.CreateProducer(); err != nil {
		return nil, err
	}

	return processor.NewProcessingService(producer, producer.Processor()), nil
}
