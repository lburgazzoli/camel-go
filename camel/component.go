package camel

// Component --
type Component interface {
	ContextAware
	Service

	CreateEndpoint(remaining string, options map[string]interface{}) (Endpoint, error)
}
