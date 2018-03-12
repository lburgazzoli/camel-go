package camel

// Component --
type Component interface {
	ContextAware

	CreateEndpoint(remaining string, options map[string]interface{}) (Endpoint, error)
}
