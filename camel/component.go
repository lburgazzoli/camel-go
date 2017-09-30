package camel

// Component --
type Component interface {
	Service
	ContextAware

	Process(message string)
}
