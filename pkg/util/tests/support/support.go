package support

import (
	"context"
	"go.uber.org/zap"
	"testing"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/stretchr/testify/assert"
)

func NewChannelVerticle(channel chan camel.Message) camel.Verticle {
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

func Run(t *testing.T, name string, fn func(*testing.T, context.Context, camel.Context)) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		l, err := zap.NewDevelopment()
		assert.Nil(t, err)

		ctx := context.Background()
		camelContext := core.NewContext(l)

		assert.NotNil(t, camelContext)

		defer func() {
			_ = camelContext.Close(ctx)
		}()

		fn(t, ctx, camelContext)
	})
}
