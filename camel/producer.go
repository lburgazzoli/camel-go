package camel

// Producer --
type Producer interface {
	Service
	Processor

	Endpoint() Endpoint
}
