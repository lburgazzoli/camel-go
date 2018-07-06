package api

// Consumer --
type Consumer interface {
	Service

	Endpoint() Endpoint
	Processor() Processor
}
