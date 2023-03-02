package engine

import (
	"context"
	"os"
	"testing"

	"github.com/lburgazzoli/camel-go/pkg/util/uuid"

	"github.com/lburgazzoli/camel-go/pkg/wasm"
	"github.com/lburgazzoli/camel-go/test/support"

	"github.com/lburgazzoli/camel-go/pkg/wasm/functions"

	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/stretchr/testify/assert"
)

func TestWASM(t *testing.T) {
	ctx := context.Background()

	r, err := wasm.NewRuntime(ctx, wasm.Options{})
	assert.Nil(t, err)

	defer func() { _ = r.Close(ctx) }()

	fd, err := os.Open("../../etc/fn/simple_process.wasm")
	assert.Nil(t, err)

	f, err := r.Load(ctx, "process", fd)
	assert.Nil(t, err)

	in, err := message.New()
	assert.Nil(t, err)

	out, err := support.Process(ctx, f, in)
	assert.Nil(t, err)

	c, ok := out.Content().([]byte)
	assert.True(t, ok)
	assert.Equal(t, "hello from wasm", string(c))

}

func TestCallbackWASM(t *testing.T) {
	ctx := context.Background()

	r, err := wasm.NewRuntime(ctx, wasm.Options{})
	assert.Nil(t, err)

	err = r.Export(ctx, "http", functions.HTTPRequest)
	assert.Nil(t, err)

	defer func() { _ = r.Close(ctx) }()

	fd, err := os.Open("../../etc/components/slack.wasm")
	assert.Nil(t, err)

	f, err := r.Load(ctx, "process", fd)
	assert.Nil(t, err)

	in, err := message.New()
	assert.Nil(t, err)
	in.SetAnnotation("webhook", "https://hooks.slack.com/services/"+uuid.New()+"/"+uuid.New()+"/"+uuid.New())
	in.SetContent("hello from gamel")

	out, err := support.Process(ctx, f, in)
	assert.Nil(t, err)

	c, ok := out.Content().([]byte)
	assert.True(t, ok)
	assert.Contains(t, string(c), "no_team")
}
