package camel

// Endpoint --
type Endpoint interface {
	Service

	Component() Component

	CreateProducer(pipe *Pipe) (Producer, error)
	CreateConsumer(pipe *Pipe) (Consumer, error)
}
