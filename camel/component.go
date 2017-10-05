package camel

// Component --
type Component interface {
	Service
	ContextAware

	Process(Exchange *Exchange)
}

// ==========================
//
// DefaultComponent
//
// ==========================

// DefaultComponent --
type DefaultComponent struct {
	context *Context
}

// SetContext --
func (component *DefaultComponent) SetContext(context *Context) {
	component.context = context
}

// Context --
func (component *DefaultComponent) Context() *Context {
	return component.context
}

// Status --
func (component *DefaultComponent) Status() ServiceStatus {
	return ServiceStatusSTARTED
}

// Process --
func (component *DefaultComponent) Process(message string) {
}
