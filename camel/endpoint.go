package camel

// Endpoint --
type Endpoint interface {
	Component() Component

	CreateProducer() (Producer, error)
	CreateConsumer() (Consumer, error)
}
