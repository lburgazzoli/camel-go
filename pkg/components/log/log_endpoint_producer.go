// //go:build components_log || components_all
package log

import (
	"context"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"go.uber.org/zap"
)

type Producer struct {
	processors.DefaultVerticle

	endpoint *Endpoint
	logger   *zap.Logger
}

func (p *Producer) Endpoint() camel.Endpoint {
	return p.endpoint
}

func (p *Producer) Start(context.Context) error {
	return nil
}

func (p *Producer) Stop(context.Context) error {
	return nil
}

func (p *Producer) Receive(ac actor.Context) {
	msg, ok := ac.Message().(camel.Message)
	if ok {
		var content string

		if _, err := p.endpoint.Component().Context().TypeConverter().Convert(msg.Content(), &content); err != nil {
			panic(err)
		}

		p.logger.Info(
			content,
			zap.String("message.type", msg.GetType()),
			zap.String("message.id", msg.GetID()))

		ac.Request(ac.Parent(), msg)
	}
}
