package choice

import (
	"context"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

type Otherwise struct {
	processors.DefaultVerticle `yaml:",inline"`

	Steps []processors.Step `yaml:"steps,omitempty"`
}

func (o *Otherwise) Configure(ctx context.Context, camelContext camel.Context) error {
	o.DefaultVerticle.SetContext(camelContext)
	return nil
}
