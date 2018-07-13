package route

import (
	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/camel"
	"github.com/rs/zerolog/log"
)

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
func (definition *FromDefinition) Unwrap(context api.Context, parent api.Processor) (api.Processor, api.Service, error) {
	var err error
	var consumer api.Consumer
	var endpoint api.Endpoint

	if endpoint, err = camel.NewEndpointFromURI(context, definition.URI); err != nil {
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
