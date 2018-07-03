package camel

import "github.com/lburgazzoli/camel-go/api"

// Producer --
type Producer interface {
	api.Service

	Endpoint() Endpoint
	Processor() Processor
}
