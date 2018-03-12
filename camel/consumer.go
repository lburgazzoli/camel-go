package camel

// Consumer --
type Consumer interface {
	Service

	Endpoint() Endpoint
}
