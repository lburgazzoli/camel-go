package core

import "github.com/lburgazzoli/camel-go/camel"

// DefaultComponent --
type DefaultComponent struct {
	DefaultService

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
