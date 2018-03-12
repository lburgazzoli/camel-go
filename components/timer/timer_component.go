package timer

import (
	"fmt"

	"github.com/lburgazzoli/camel-go/camel"
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
	endpoint := timerEndpoint{}
	endpoint.component = component

	// endpoint option validation
	if _, ok := options["duration"]; !ok {
		return nil, fmt.Errorf("Missing mandatory option: duration")
	}

	// bind options to endpoint
	camel.SetFields(component.context, &endpoint, options)

	return &endpoint, nil
}
