// //go:build steps_transform || steps_all

package transform

import (
	"context"
	"testing"
	"time"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
	"github.com/stretchr/testify/assert"
)

func TestTransformJQ(t *testing.T) {
	support.Run(t, "jq", func(t *testing.T, ctx context.Context, c camel.Context) {
		t.Helper()
		t.Skip("TODO")

		wg := make(chan camel.Message)

		v := support.NewChannelVerticle(wg)
		err := c.Spawn(v)

		assert.Nil(t, err)

		p := Transform{
			DefaultVerticle: processors.NewDefaultVerticle(),
			Language:        Language{}}

		p.Next(v.ID())

		id, err := p.Reify(ctx, c)
		assert.Nil(t, err)
		assert.NotNil(t, id)

		msg, err := message.New()
		assert.Nil(t, err)

		msg.SetAnnotation("foo", "bar")

		assert.Nil(t, c.Send(id, msg))

		select {
		case msg := <-wg:
			c, ok := msg.Content().([]byte)
			assert.True(t, ok)
			assert.Equal(t, "hello "+msg.GetID()+", bar", string(c))
		case <-time.After(5 * time.Second):
			assert.Fail(t, "timeout")
		}
	})
}
