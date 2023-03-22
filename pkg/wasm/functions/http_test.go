package functions

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/lburgazzoli/camel-go/pkg/wasm"
	"github.com/stretchr/testify/require"

	"github.com/lburgazzoli/camel-go/pkg/util/uuid"

	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/stretchr/testify/assert"
)

func TestHttpFunction(t *testing.T) {
	ctx := context.Background()

	r, err := wasm.NewRuntime(ctx, wasm.Options{})
	assert.Nil(t, err)

	err = r.Export(ctx, "http", HTTPRequest)
	assert.Nil(t, err)

	defer func() { _ = r.Close(ctx) }()

	data, err := filepath.Abs("../../../etc/wasm/components/slack.wasm")
	require.NoError(t, err)
	require.FileExists(t, data)

	fd, err := os.Open(data)
	assert.Nil(t, err)

	f, err := r.Load(ctx, "process", fd)
	assert.Nil(t, err)

	in, err := message.New()
	assert.Nil(t, err)
	in.SetAnnotation("webhook", "https://hooks.slack.com/services/"+uuid.New()+"/"+uuid.New()+"/"+uuid.New())
	in.SetContent("hello from gamel")

	out, err := wasm.Process(ctx, f, in)
	assert.Nil(t, err)

	c, ok := out.Content().([]byte)
	assert.True(t, ok)
	assert.Contains(t, string(c), "no_team")
}
