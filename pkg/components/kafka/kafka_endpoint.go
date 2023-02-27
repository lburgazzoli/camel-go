package kafka

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

type Endpoint struct {
	config Config
	components.DefaultEndpoint
}

func (e *Endpoint) Start() error {
	return nil
}

func (e *Endpoint) Stop() error {
	return nil
}

func (e *Endpoint) Producer() (api.Producer, error) {
	c := Producer{
		id:       uuid.New(),
		endpoint: e,
	}

	return &c, nil
}
