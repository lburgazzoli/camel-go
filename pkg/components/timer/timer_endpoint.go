////go:build components_timer || components_all

package timer

import (
	"context"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
)

type Endpoint struct {
	components.DefaultEndpoint

	config Config
}

func (e *Endpoint) Start(context.Context) error {
	return nil
}

func (e *Endpoint) Stop(context.Context) error {
	return nil
}

func (e *Endpoint) Consumer() (api.Consumer, error) {
	c := Consumer{
		DefaultVerticle: processors.NewDefaultVerticle(),
		endpoint:        e,
	}

	return &c, nil
}
