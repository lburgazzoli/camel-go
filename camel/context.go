package camel

// Context --
type Context interface {
	Service

	AddRegistryLoader(loader RegistryLoader)

	AddComponent(name string, component Component)

	GetComponent(name string) (Component, error)
}
