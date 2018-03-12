package log

import (
	"github.com/lburgazzoli/camel-go/camel"
	"github.com/rs/zerolog"
)

// ==========================
//
//
//
// ==========================

// NewComponent --
func NewComponent() camel.Component {
	return &Component{}
}

// ==========================
//
// Component
//
// ==========================

// Component --
type Component struct {
	context *camel.Context
}

// SetContext --
func (component *Component) SetContext(context *camel.Context) {
	component.context = context
}

// Context --
func (component *Component) Context() *camel.Context {
	return component.context
}

// CreateEndpoint --
func (component *Component) CreateEndpoint(remaining string, options map[string]interface{}) (camel.Endpoint, error) {
	// Create the endpoint and set default values
	endpoint := logEndpoint{component: component}
	endpoint.name = remaining
	endpoint.level = zerolog.InfoLevel

	// bind options to endpoint
	camel.SetFields(component.context, &endpoint, options)

	return &endpoint, nil
}
