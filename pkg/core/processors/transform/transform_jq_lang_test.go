// //go:build steps_transform || steps_all

package transform

import (
	"context"
	"testing"
	"time"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/language"
	"github.com/lburgazzoli/camel-go/pkg/core/language/jq"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
)

func TestTransformJQ(t *testing.T) {
	support.Run(t, "jq", func(t *testing.T, ctx context.Context) {
		t.Helper()

		wg := make(chan camel.Message)

		c := camel.GetContext(ctx)

		wgv, err := support.NewChannelVerticle(wg).Reify(ctx)
		require.Nil(t, err)
		require.NotNil(t, wgv)

		wgp, err := c.Spawn(wgv)
		require.Nil(t, err)
		require.NotNil(t, wgp)

		l := language.Language{
			Jq: &jq.Jq{
				Expression: `.message`,
			},
		}

		pv, err := NewTransformWithLanguage(l).Reify(ctx)
		require.Nil(t, err)
		require.NotNil(t, pv)

		pv.Next(wgp)

		pvp, err := c.Spawn(pv)
		require.Nil(t, err)
		require.NotNil(t, pvp)

		msg, err := message.New()
		require.Nil(t, err)

		msg.SetContent(`{ "message": "hello jq" }`)
		msg.SetAnnotation("foo", "bar")

		require.Nil(t, c.SendTo(pvp, msg))

		select {
		case msg := <-wg:
			c, ok := msg.Content().(string)
			assert.True(t, ok)
			assert.Equal(t, "hello jq", c)
		case <-time.After(5 * time.Second):
			assert.Fail(t, "timeout")
		}
	})
}
