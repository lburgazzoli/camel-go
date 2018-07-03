package camel

import "github.com/lburgazzoli/camel-go/api"

// Consumer --
type Consumer interface {
	api.Service

	Endpoint() Endpoint
	Processor() Processor
}
