package wasm

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	camel "github.com/lburgazzoli/camel-go/pkg/api"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestWASM(t *testing.T) {
	support.Run(t, "wasm", func(t *testing.T, ctx context.Context) {
		t.Helper()

		r, err := NewRuntime(ctx)
		assert.Nil(t, err)

		c := camel.ExtractContext(ctx)

		defer func() { _ = r.Close(ctx) }()

		path, err := filepath.Abs("../../etc/wasm/fn/simple_process.wasm")
		require.NoError(t, err)
		require.FileExists(t, path)

		content, err := os.Open(path)
		require.NoError(t, err)

		m, err := r.Load(ctx, content)
		require.NoError(t, err)

		defer func() { _ = m.Close(ctx) }()

		p, err := m.Processor(ctx)
		require.NoError(t, err)

		in := c.NewMessage()

		err = p.Process(ctx, in)
		require.NoError(t, err)

		data, ok := in.Content().([]byte)
		assert.True(t, ok)
		assert.Equal(t, "hello from wasm", string(data))

	})
}
