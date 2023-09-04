// //go:build components_log || components_all
package log

import (
	"context"
	"log/slog"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
)

type Producer struct {
	processors.DefaultVerticle

	endpoint *Endpoint
	logger   *slog.Logger
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
			slog.Group(
				"message",
				slog.String("id", msg.ID()),
				slog.String("type", msg.Type()),
			),
		)

		ac.Request(ac.Parent(), msg)
	}
}
