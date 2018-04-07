package timer

import (
	"fmt"

	"github.com/lburgazzoli/camel-go/camel"
	"github.com/lburgazzoli/camel-go/introspection"
)

// ==========================
//
//
//
// ==========================

// NewComponent --
func NewComponent() camel.Component {
	return &Component{
		state: camel.NewServiceState(camel.ServiceStatusSTOPPED),
	}
}

// ==========================
//
// Component
//
// ==========================

// Component --
type Component struct {
	state   camel.ServiceState
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

// Start --
func (component *Component) Start() {
	component.state.Transition(camel.ServiceStatusSTOPPED, camel.ServiceStatusSTARTED, component.doStart)
	component.state.Transition(camel.ServiceStatusSUSPENDED, camel.ServiceStatusSTARTED, component.doStart)
}

// Stop --
func (component *Component) Stop() {
	component.state.Transition(camel.ServiceStatusSTARTED, camel.ServiceStatusSTOPPED, component.doStop)
}

// CreateEndpoint --
func (component *Component) CreateEndpoint(remaining string, options map[string]interface{}) (camel.Endpoint, error) {
	// Create the endpoint and set default values
	endpoint := timerEndpoint{}
	endpoint.component = component

	// endpoint option validation
	if _, ok := options["period"]; !ok {
		return nil, fmt.Errorf("missing mandatory option: period")
	}

	// bind options to endpoint
	introspection.SetProperties(component.context, &endpoint, options)

	return &endpoint, nil
}

// ==========================
//
// Helpers
//
// ==========================

func (component *Component) doStart() {
}

func (component *Component) doStop() {
}
