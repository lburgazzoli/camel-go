// //go:build steps_transform || steps_all

package transform

import (
	"context"
	"testing"
	"time"

	"github.com/lburgazzoli/camel-go/pkg/core/language/mustache"

	"github.com/lburgazzoli/camel-go/pkg/core/language"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
	"github.com/stretchr/testify/assert"
)

func TestTransformMustache(t *testing.T) {
	support.Run(t, "mustache", func(t *testing.T, ctx context.Context) {
		t.Helper()

		c := camel.GetContext(ctx)

		wg := make(chan camel.Message)

		v := support.NewChannelVerticle(wg)
		err := c.Spawn(v)

		assert.Nil(t, err)

		p := Transform{
			DefaultVerticle: processors.NewDefaultVerticle(),
			Language: language.Language{
				Mustache: &mustache.Mustache{
					Template: `hello {{message.id}}, {{message.annotations.foo}}`,
				},
			}}

		p.Next(v.ID())

		id, err := p.Reify(ctx)
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
