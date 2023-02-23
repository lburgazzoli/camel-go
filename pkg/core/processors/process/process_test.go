package process

import (
	"context"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"

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

	p := Process{Ref: "p"}
	p.Next(c.SpawnFn(func(c actor.Context) {
		switch msg := c.Message().(type) {
		case api.Message:
			wg <- msg
		}
	}))

	pid, err := p.Reify(c)
	assert.Nil(t, err)
	assert.NotNil(t, pid)

	msg, err := message.New()
	assert.Nil(t, err)

	c.Send(pid, msg)

	select {
	case msg := <-wg:
		assert.Equal(t, content, msg.Content())
	case <-time.After(5 * time.Second):
		assert.Fail(t, "timeout")
	}
}
