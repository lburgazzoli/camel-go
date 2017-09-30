package camel

// Context --
type Context interface {
	Service

	AddRegistryLoader(loader RegistryLoader)

	AddComponent(name string, component Component)

	Component(name string) (Component, error)
}

// ContextAware --
type ContextAware interface {
	SetContext(context Context)
	Context() Context
}
