// //go:build components_log || components_all
package log

import (
	"context"
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
)

type Producer struct {
	camel.WithOutputs

	id       string
	endpoint *Endpoint
}

func (p *Producer) ID() string {
	return p.id
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

func (p *Producer) Receive(ctx actor.Context) {
	msg, ok := ctx.Message().(camel.Message)
	if ok {
		fmt.Printf("(%s) %s/%s -> %s\n", p.endpoint.config.Remaining, msg.GetType(), msg.GetID(), msg.Content())
	}
}
