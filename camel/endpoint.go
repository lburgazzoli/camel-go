package camel

// Endpoint --
type Endpoint interface {
	Component() Component
	URI() string

	NewProducer() (Producer, error)
	NewConsumer() (Consumer, error)
}
