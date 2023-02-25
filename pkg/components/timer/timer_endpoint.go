////go:build components_timer || components_all

package timer

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
)

type Endpoint struct {
	components.DefaultEndpoint

	config Config
}

func (e *Endpoint) Start() error {
	return nil
}

func (e *Endpoint) Stop() error {
	return nil
}

func (e *Endpoint) Consumer() (api.Consumer, error) {
	c := Consumer{
		endpoint: e,
	}

	return &c, nil
}
