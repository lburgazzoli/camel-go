// //go:build steps_process || steps_all

package setbody

import (
	"context"
	"testing"
	"time"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/stretchr/testify/require"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
)

func TestSetBody(t *testing.T) {
	support.Run(t, "simple", func(t *testing.T, ctx context.Context) {
		t.Helper()

		content := uuid.New()

		c := camel.ExtractContext(ctx)
		c.Registry().Set("p", func(_ context.Context, message camel.Message) error {
			message.SetContent(content)
			return nil
		})

		p := New()
		p.Language.Constant = &LanguageConstant{
			Value: content,
		}

		pv, err := p.Reify(ctx)
		require.Nil(t, err)
		require.NotNil(t, pv)

		pvp, err := c.Spawn(pv)
		require.Nil(t, err)
		require.NotNil(t, pvp)

		msg, err := message.New()
		require.Nil(t, err)

		res, err := c.RequestTo(pvp, msg, 1*time.Second)
		require.Nil(t, err)
		require.Equal(t, content, res.Content())
	})
}
