package camel

// Producer --
type Producer interface {
	Processor

	Endpoint() Endpoint
}
