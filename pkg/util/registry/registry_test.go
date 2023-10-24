package registry

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPull(t *testing.T) {
	ctx := context.Background()

	root, err := Pull(ctx, "quay.io/lburgazzoli/camel-go-wasm:latest")

	defer func() {
		if root != "" {
			_ = os.RemoveAll(root)
		}
	}()

	require.NoError(t, err)
	assert.NotEmpty(t, root)

	fi, err := os.Stat(path.Join(root, "etc/wasm/fn/simple_process.wasm"))

	require.NoError(t, err)
	assert.False(t, fi.IsDir())
	assert.True(t, fi.Mode().IsRegular())
}

func TestBlob(t *testing.T) {
	ctx := context.Background()

	content, err := Blob(
		ctx,
		"quay.io/lburgazzoli/camel-go-wasm:latest",
		"etc/wasm/fn/simple_process.wasm")

	defer func() {
		if content != nil {
			_ = content.Close()
		}
	}()

	require.NoError(t, err)
	assert.NotNil(t, content)
}
