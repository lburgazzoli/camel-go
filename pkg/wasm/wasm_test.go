package wasm

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/stretchr/testify/assert"
)

func TestWASM(t *testing.T) {
	ctx := context.Background()

	r, err := NewRuntime(ctx, Options{})
	assert.Nil(t, err)

	defer func() { _ = r.Close(ctx) }()

	data, err := filepath.Abs("../../etc/wasm/fn/simple_process.wasm")
	require.NoError(t, err)
	require.FileExists(t, data)

	f, err := r.Load(ctx, data)
	assert.Nil(t, err)

	in, err := message.New()
	assert.Nil(t, err)

	out, err := f.Invoke(ctx, in)
	assert.Nil(t, err)

	c, ok := out.Content().([]byte)
	assert.True(t, ok)
	assert.Equal(t, "hello from wasm", string(c))

}
