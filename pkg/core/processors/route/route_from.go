package route

import (
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/endpoint"
)

type From struct {
	endpoint.Endpoint `yaml:",inline"`

	Steps []processors.Step `yaml:"steps,omitempty"`
}
