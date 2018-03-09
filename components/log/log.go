package log

import (
	golog "log"

	"github.com/lburgazzoli/camel-go/camel"
)

// ==========================
//
//
//
// ==========================

// Component --
type Component struct {
	context *camel.Context
	logger  *golog.Logger
}

// SetContext --
func (component *Component) SetContext(context *camel.Context) {
	component.context = context
}

// Context --
func (component *Component) Context() *camel.Context {
	return component.context
}

// Process --
func (component *Component) Process(exchange *camel.Exchange) {
	component.logger.Printf("%+v\n", exchange.Body())
}

// ==========================
//
//
//
// ==========================

// NewComponent --
func NewComponent() camel.Component {
	return &Component{
		logger: new(golog.Logger),
	}
}
