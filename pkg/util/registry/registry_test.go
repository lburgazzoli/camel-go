package registry

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrasCopyFS(t *testing.T) {
	ctx := context.Background()

	root, err := Pull(ctx, "docker.io/lburgazzoli/camel-go:latest")

	defer func() {
		if root != "" {
			_ = os.RemoveAll(root)
		}
	}()

	assert.Nil(t, err)
	assert.NotEmpty(t, root)

	fi, err := os.Stat(path.Join(root, "etc/fn/simple_process.wasm"))

	assert.Nil(t, err)
	assert.False(t, fi.IsDir())
	assert.True(t, fi.Mode().IsRegular())
}
