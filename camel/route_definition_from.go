package camel

import "github.com/rs/zerolog/log"

// FromDefinition --
type FromDefinition struct {
	parent   Definition
	children []Definition

	URI string
}

// Parent --
func (definition *FromDefinition) Parent() Definition {
	return definition.parent
}

// Children --
func (definition *FromDefinition) Children() []Definition {
	return definition.children
}

// Unwrap ---
func (definition *FromDefinition) Unwrap(context *Context, parent Processor) (Processor, Service, error) {
	var err error
	var consumer Consumer
	var endpoint Endpoint

	if endpoint, err = context.CreateEndpointFromURI(definition.URI); err != nil {
		return parent, nil, nil
	}

	if consumer, err = endpoint.CreateConsumer(); err != nil {
		return parent, nil, nil
	}

	if parent != nil {
		log.Panic().Msgf("parent pipe should be nil, got %+v", parent)
	}

	return consumer.Processor(), consumer, nil
}
