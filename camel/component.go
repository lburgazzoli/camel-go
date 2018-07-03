package camel

import "github.com/lburgazzoli/camel-go/api"

// Component --
type Component interface {
	ContextAware
	api.Service

	CreateEndpoint(remaining string, options map[string]interface{}) (Endpoint, error)
}
