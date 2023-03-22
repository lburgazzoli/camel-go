// //go:build steps_process || steps_all

package process

import (
	"context"
	"testing"
	"time"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
)

func TestProcessor(t *testing.T) {
	support.Run(t, "simple", func(t *testing.T, ctx context.Context) {
		t.Helper()

		content := uuid.New()
		wg := make(chan camel.Message)

		c := camel.GetContext(ctx)
		c.Registry().Set("p", func(_ context.Context, message camel.Message) error {
			message.SetContent(content)
			return nil
		})

		wgv, err := support.NewChannelVerticle(wg).Reify(ctx)
		require.Nil(t, err)
		require.NotNil(t, wgv)

		wgp, err := c.Spawn(wgv)
		require.Nil(t, err)
		require.NotNil(t, wgp)

		pv, err := NewProcessWithRef("p").Reify(ctx)
		require.Nil(t, err)
		require.NotNil(t, pv)

		pv.Next(wgp)

		pvp, err := c.Spawn(pv)
		require.Nil(t, err)
		require.NotNil(t, pvp)

		msg, err := message.New()
		require.Nil(t, err)
		require.Nil(t, c.SendTo(pvp, msg))

		select {
		case msg := <-wg:
			assert.Equal(t, content, msg.Content())
		case <-time.After(5 * time.Second):
			assert.Fail(t, "timeout")
		}
	})
}
