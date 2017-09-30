package camel

// Registry --
type Registry interface { 
}

// RegistryLoader --
type RegistryLoader interface {
	Service

	Load(name string)  (interface{}, error) 
}