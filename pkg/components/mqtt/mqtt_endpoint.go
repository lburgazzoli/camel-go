////go:build components_mqtt || components_all

package mqtt

import (
	"context"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"

	"go.uber.org/zap"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
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
	id := uuid.New()

	c := Consumer{
		DefaultVerticle: processors.NewDefaultVerticle(),
		endpoint:        e,
		logger:          e.Logger().With(zap.String("consumer.id", id)).Sugar(),
	}

	return &c, nil
}
