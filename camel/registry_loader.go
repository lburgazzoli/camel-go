package camel

// RegistryLoader --
type RegistryLoader interface {
	Service

	Load(name string) (interface{}, error)
}
