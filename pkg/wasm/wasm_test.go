package wasm

import (
	"context"
	"testing"

	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/stretchr/testify/assert"
)

func TestWASM(t *testing.T) {
	ctx := context.Background()

	r, err := NewRuntime(ctx, Options{})
	assert.Nil(t, err)

	defer func() { _ = r.Close(ctx) }()

	p, err := r.Load(ctx, "../../etc/fn/simple_process.wasm")
	assert.Nil(t, err)

	in, err := message.New()
	assert.Nil(t, err)

	out, err := p.Process(ctx, in)
	assert.Nil(t, err)

	c, ok := out.Content().([]byte)
	assert.True(t, ok)
	assert.Equal(t, "hello from wasm", string(c))

}
