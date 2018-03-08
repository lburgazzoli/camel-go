package camel

// HasContext --
type HasContext interface {
	Context() *Context
}

// ContextAware --
type ContextAware interface {
	HasContext

	SetContext(context *Context)
}
