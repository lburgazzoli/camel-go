package core

import "github.com/lburgazzoli/camel-go/camel"

// NewComponent --
func NewComponent(order int) camel.Component {
	return &DefaultComponent{
		Service: NewServiceWithOrder(order),
	}
}

// DefaultComponent --
type DefaultComponent struct {
	camel.Service

	context camel.Context
}

// Process --
func (component *DefaultComponent) Process(message string) {
}

// SetContext --
func (component *DefaultComponent) SetContext(context camel.Context) {
	component.context = context
}

// Context --
func (component *DefaultComponent) Context() camel.Context {
	return component.context
}
