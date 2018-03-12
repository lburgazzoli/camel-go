package camel

// Endpoint --
type Endpoint interface {
	Service

	Component() Component

	CreateProducer() (Producer, error)
	CreateConsumer(processor Processor) (Consumer, error)
}
