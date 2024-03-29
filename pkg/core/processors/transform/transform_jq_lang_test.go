// //go:build steps_transform || steps_all

package transform

import (
	"testing"
	"time"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/language"
	"github.com/lburgazzoli/camel-go/pkg/core/language/jq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
)

func TestTransformJQ(t *testing.T) {
	g := support.With(t)
	c := camel.ExtractContext(g.Ctx())

	wg := make(chan camel.Message)

	wgv, err := support.NewChannelVerticle(wg).Reify(g.Ctx())
	require.NoError(t, err)
	require.NotNil(t, wgv)

	wgp, err := c.Spawn(wgv)
	require.NoError(t, err)
	require.NotNil(t, wgp)

	l := language.Language{
		Jq: &jq.Jq{
			Definition: jq.Definition{Expression: `.message`},
		},
	}

	pv, err := New(WithLanguage(l)).Reify(g.Ctx())
	require.NoError(t, err)
	require.NotNil(t, pv)

	pvp, err := c.Spawn(pv)
	require.NoError(t, err)
	require.NotNil(t, pvp)

	msg := c.NewMessage()
	msg.SetContent(`{ "message": "hello jq" }`)
	msg.SetAttribute("foo", "bar")

	res, err := c.RequestTo(pvp, msg, 1*time.Second)
	require.NoError(t, err)

	body, ok := res.Content().(string)
	assert.True(t, ok)
	assert.Equal(t, "hello jq", body)
}
