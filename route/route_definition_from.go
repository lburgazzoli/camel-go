package route

import (
	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/camel"
	"github.com/lburgazzoli/camel-go/processor"
)

// FromDefinition --
type FromDefinition struct {
	api.ContextAware
	ProcessingNode

	context  api.Context
	parent   Definition
	children []Definition

	URI string
}

// SetContext --
func (definition *FromDefinition) SetContext(context api.Context) {
	definition.context = context
}

// Context --
func (definition *FromDefinition) Context() api.Context {
	return definition.context
}

// Parent --
func (definition *FromDefinition) Parent() Definition {
	return definition.parent
}

// Children --
func (definition *FromDefinition) Children() []Definition {
	return definition.children
}

// Processor ---
func (definition *FromDefinition) Processor() (api.Processor, error) {
	var err error
	var consumer api.Consumer
	var endpoint api.Endpoint

	if endpoint, err = camel.NewEndpointFromURI(definition.context, definition.URI); err != nil {
		return nil, nil
	}

	if consumer, err = endpoint.CreateConsumer(); err != nil {
		return nil, nil
	}

	return processor.NewProcessingService(consumer, consumer.Processor()), nil
}
