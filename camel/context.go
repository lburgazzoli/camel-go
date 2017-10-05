package camel

// TypeConverter --
type TypeConverter interface {
}

// Context --
type Context interface {
	Service

	AddRegistryLoader(loader RegistryLoader)

	AddTypeConverter(converter TypeConverter)

	AddComponent(name string, component Component)

	Component(name string) (Component, error)
}

// ContextAware --
type ContextAware interface {
	SetContext(context Context)
	Context() Context
}
