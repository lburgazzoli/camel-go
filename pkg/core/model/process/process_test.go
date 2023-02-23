package process

import (
	"context"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimpleProcessor(t *testing.T) {
	content := uuid.New()
	wg := make(chan api.Message)

	c := core.NewContext(context.Background())
	assert.NotNil(t, c)

	c.Registry().Set("p", func(message api.Message) {
		message.SetContent(content)
	})

	as := actor.NewActorSystem()
	receiverPid := as.Root.Spawn(actor.PropsFromFunc(func(c actor.Context) {
		switch msg := c.Message().(type) {
		case *actor.Started:
			break
		case *actor.Stopping:
			break
		case *actor.Stopped:
			break
		case *actor.Restarting:
			break
		case api.Message:
			wg <- msg
		}
	}))

	p := Process{
		Ref: "p",
		SendTo: []*actor.PID{
			receiverPid,
		},
	}

	a, err := p.Reify(c)
	assert.Nil(t, err)

	senderPid := as.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return a
	}))

	msg, err := message.New()
	assert.Nil(t, err)

	as.Root.Send(senderPid, msg)

	select {
	case msg := <-wg:
		assert.Equal(t, content, msg.Content())
	case <-time.After(5 * time.Second):
		assert.Fail(t, "timeout")
	}
}
