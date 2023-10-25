package wasm

import (
	"os"
	"path/filepath"
	"testing"

	camel "github.com/lburgazzoli/camel-go/pkg/api"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestWASM(t *testing.T) {
	g := support.With(t)
	c := camel.ExtractContext(g.Ctx())

	r, err := NewRuntime(g.Ctx())
	require.NoError(t, err)

	defer func() { _ = r.Close(g.Ctx()) }()

	path, err := filepath.Abs("../../etc/wasm/fn/simple_process.wasm")
	require.NoError(t, err)
	require.FileExists(t, path)

	content, err := os.Open(path)
	require.NoError(t, err)

	m, err := r.Load(g.Ctx(), content)
	require.NoError(t, err)

	defer func() { _ = m.Close(g.Ctx()) }()

	p, err := m.Processor(g.Ctx())
	require.NoError(t, err)

	in := c.NewMessage()

	err = p.Process(g.Ctx(), in)
	require.NoError(t, err)

	data, ok := in.Content().([]byte)
	assert.True(t, ok)
	assert.Equal(t, "hello from wasm", string(data))
}
