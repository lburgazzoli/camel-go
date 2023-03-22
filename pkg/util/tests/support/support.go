package support

import (
	"context"
	"testing"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"

	"go.uber.org/zap"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/stretchr/testify/assert"
)

type ReifyableVerticle interface {
	camel.Verticle
	processors.Reifyable
}

func NewChannelVerticle(channel chan camel.Message) ReifyableVerticle {
	return &ChannelVerticle{
		id:      uuid.New(),
		channel: channel,
	}
}

type ChannelVerticle struct {
	camel.WithOutputs

	id      string
	channel chan camel.Message
}

func (p *ChannelVerticle) ID() string {
	return p.id
}

func (p *ChannelVerticle) Receive(c actor.Context) {
	msg, ok := c.Message().(camel.Message)
	if ok {
		p.channel <- msg
	}
}

func (p *ChannelVerticle) Reify(_ context.Context) (camel.Verticle, error) {
	return p, nil
}

func Run(t *testing.T, name string, fn func(*testing.T, context.Context)) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		l, err := zap.NewDevelopment()
		assert.Nil(t, err)

		camelContext := core.NewContext(l)
		ctx := context.WithValue(context.Background(), camel.ContextKeyCamelContext, camelContext)

		assert.NotNil(t, camelContext)

		defer func() {
			_ = camelContext.Close(ctx)
		}()

		fn(t, ctx)
	})
}
