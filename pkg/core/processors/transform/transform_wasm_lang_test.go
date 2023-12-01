// //go:build steps_transform || steps_all

package transform

import (
	"testing"
	"time"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/language"
	"github.com/lburgazzoli/camel-go/pkg/core/language/wasm"
	"github.com/lburgazzoli/camel-go/pkg/util/registry"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
)

func TestTransformWASMLocal(t *testing.T) {

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
		Wasm: wasm.New(wasm.WithRef("../../../../etc/wasm/fn/simple_process.wasm")),
	}

	pv, err := New(WithLanguage(l)).Reify(g.Ctx())
	require.NoError(t, err)
	require.NotNil(t, pv)

	pvp, err := c.Spawn(pv)
	require.NoError(t, err)
	require.NotNil(t, pvp)

	msg := c.NewMessage()

	msg.SetContent(uuid.New())
	msg.SetAttribute("foo", "bar")

	res, err := c.RequestTo(pvp, msg, 1*time.Second)
	require.NoError(t, err)

	body, ok := res.Content().([]byte)
	assert.True(t, ok)
	assert.Equal(t, "hello from wasm", string(body))
}

func TestTransformWASMRegistry(t *testing.T) {
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
		Wasm: wasm.New(wasm.WithRef(registry.WasmContainerImage() + "?etc/wasm/fn/simple_process.wasm")),
	}

	pv, err := New(WithLanguage(l)).Reify(g.Ctx())
	require.NoError(t, err)
	require.NotNil(t, pv)

	pvp, err := c.Spawn(pv)
	require.NoError(t, err)
	require.NotNil(t, pvp)

	msg := c.NewMessage()

	msg.SetContent(uuid.New())
	msg.SetAttribute("foo", "bar")

	res, err := c.RequestTo(pvp, msg, 1*time.Second)
	require.NoError(t, err)

	body, ok := res.Content().([]byte)
	assert.True(t, ok)
	assert.Equal(t, "hello from wasm", string(body))

	require.NoError(t, c.SendTo(pvp, msg))
}
