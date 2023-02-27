package support

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

func NewChannelVerticle(channel chan api.Message) api.Verticle {
	return &ChannelVerticle{
		id:      uuid.New(),
		channel: channel,
	}
}

type ChannelVerticle struct {
	api.WithOutputs

	id      string
	channel chan api.Message
}

func (p *ChannelVerticle) ID() string {
	return p.id
}

func (p *ChannelVerticle) Receive(c actor.Context) {
	msg, ok := c.Message().(api.Message)
	if ok {
		p.channel <- msg
	}
}
