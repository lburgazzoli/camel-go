package camel

// Producer --
type Producer interface {
	Endpoint() Endpoint
}
