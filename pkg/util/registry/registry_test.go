package registry

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPull(t *testing.T) {
	ctx := context.Background()

	root, err := Pull(ctx, "docker.io/lburgazzoli/camel-go:latest")

	defer func() {
		if root != "" {
			_ = os.RemoveAll(root)
		}
	}()

	assert.Nil(t, err)
	assert.NotEmpty(t, root)

	fi, err := os.Stat(path.Join(root, "etc/wasm/fn/simple_process.wasm"))

	assert.Nil(t, err)
	assert.False(t, fi.IsDir())
	assert.True(t, fi.Mode().IsRegular())
}

func TestBlob(t *testing.T) {
	ctx := context.Background()

	content, err := Blob(
		ctx,
		"docker.io/lburgazzoli/camel-go:latest",
		"etc/wasm/fn/simple_process.wasm")

	defer func() {
		if content != nil {
			_ = content.Close()
		}
	}()

	assert.Nil(t, err)
	assert.NotNil(t, content)
}
