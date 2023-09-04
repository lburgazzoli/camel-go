// //go:build components_dapr_pubsub || components_all

package pubsub

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
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

func (e *Endpoint) Consumer(pid *actor.PID) (api.Consumer, error) {
	c := Consumer{
		DefaultConsumer: components.NewDefaultConsumer(e, pid),
		endpoint:        e,
	}

	return &c, nil
}

func (e *Endpoint) Producer() (api.Producer, error) {
	c := Producer{
		DefaultVerticle: processors.NewDefaultVerticle(),
		endpoint:        e,
		tc:              e.Context().TypeConverter(),
	}

	return &c, nil
}
