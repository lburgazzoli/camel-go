package camel

// ContextAware --
type ContextAware interface {
	SetContext(context Context)
	GetContext() Context
}
