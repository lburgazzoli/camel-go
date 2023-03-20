// //go:build steps_process || steps_all

package process

import (
	"context"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"testing"
	"time"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"

	"github.com/stretchr/testify/assert"
)

func TestProcessor(t *testing.T) {
	support.Run(t, "simple", func(t *testing.T, ctx context.Context, c camel.Context) {
		t.Helper()

		content := uuid.New()
		wg := make(chan camel.Message)

		c.Registry().Set("p", func(message camel.Message) {
			message.SetContent(content)
		})

		v := support.NewChannelVerticle(wg)
		err := c.Spawn(v)

		assert.Nil(t, err)

		p := Process{
			DefaultVerticle: processors.NewDefaultVerticle(),
			Ref:             "p",
		}

		p.Next(v.ID())

		id, err := p.Reify(ctx, c)
		assert.Nil(t, err)
		assert.NotNil(t, id)

		msg, err := message.New()
		assert.Nil(t, err)

		assert.Nil(t, c.Send(id, msg))

		select {
		case msg := <-wg:
			assert.Equal(t, content, msg.Content())
		case <-time.After(5 * time.Second):
			assert.Fail(t, "timeout")
		}
	})
}
