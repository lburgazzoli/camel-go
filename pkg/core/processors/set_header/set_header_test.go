// //go:build steps_process || steps_all

package setheader

import (
	"context"
	"testing"
	"time"

	"github.com/lburgazzoli/camel-go/pkg/core/language"
	"github.com/lburgazzoli/camel-go/pkg/core/language/constant"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/stretchr/testify/require"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
)

func TestSetHeader(t *testing.T) {
	support.Run(t, "constant", func(t *testing.T, ctx context.Context) {
		t.Helper()

		name := uuid.New()
		content := uuid.New()

		c := camel.ExtractContext(ctx)

		p := New(
			WithName(name),
			WithLanguage(language.Language{
				Constant: &constant.Constant{
					Value: content,
				},
			}))

		pv, err := p.Reify(ctx)
		require.Nil(t, err)
		require.NotNil(t, pv)

		pvp, err := c.Spawn(pv)
		require.Nil(t, err)
		require.NotNil(t, pvp)

		msg := c.NewMessage()

		res, err := c.RequestTo(pvp, msg, 1*time.Second)
		require.Nil(t, err)

		h, ok := res.Header(name)
		require.True(t, ok)
		require.Equal(t, content, h)
	})

	support.Run(t, "jq", func(t *testing.T, ctx context.Context) {
		t.Helper()

		name := uuid.New()

		c := camel.ExtractContext(ctx)

		p := New(
			WithName(name),
			WithLanguage(*language.New(
				language.WithJqExpression(".foo")),
			),
		)

		pv, err := p.Reify(ctx)
		require.Nil(t, err)
		require.NotNil(t, pv)

		pvp, err := c.Spawn(pv)
		require.Nil(t, err)
		require.NotNil(t, pvp)

		msg := c.NewMessage()
		msg.SetContent(`{ "foo": "bar"}`)

		res, err := c.RequestTo(pvp, msg, 1*time.Second)
		require.Nil(t, err)

		h, ok := res.Header(name)
		require.True(t, ok)
		require.Equal(t, "bar", h)
	})
}
