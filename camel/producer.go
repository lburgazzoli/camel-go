package camel

// Producer --
type Producer interface {
	Service

	Endpoint() Endpoint
	Pipe() *Pipe
}
