package api

// Producer --
type Producer interface {
	Service

	Endpoint() Endpoint
	Processor() Processor
}
