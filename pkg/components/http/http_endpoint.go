////go:build components_http || components_all

package http

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
	"github.com/lburgazzoli/camel-go/pkg/core/errors"
)

type Endpoint struct {
	config Config
	components.DefaultEndpoint
}

func (e *Endpoint) Start(context.Context) error {
	return nil
}

func (e *Endpoint) Stop(context.Context) error {
	return nil
}

func (e *Endpoint) Producer() (api.Producer, error) {
	c := Producer{
		DefaultProducer: components.NewDefaultProducer(e),
		endpoint:        e,
		tc:              e.Context().TypeConverter(),
	}

	return &c, nil
}

func (e *Endpoint) Consumer(_ *actor.PID) (api.Consumer, error) {
	return nil, errors.NotImplementedf("TODO")
}
