package camel

import "github.com/lburgazzoli/camel-go/api"

// Endpoint --
type Endpoint interface {
	api.Service

	Component() Component

	CreateProducer() (Producer, error)
	CreateConsumer() (Consumer, error)
}
