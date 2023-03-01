package wasm

import (
	"context"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/wasm/serdes"
	"os"
	"testing"

	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/stretchr/testify/assert"
)

func TestWASM(t *testing.T) {
	ctx := context.Background()

	r, err := NewRuntime(ctx, Options{})
	assert.Nil(t, err)

	defer func() { _ = r.Close(ctx) }()

	fd, err := os.Open("../../etc/fn/simple_process.wasm")
	assert.Nil(t, err)

	f, err := r.Load(ctx, "process", fd)
	assert.Nil(t, err)

	in, err := message.New()
	assert.Nil(t, err)

	out, err := process(ctx, f, in)
	assert.Nil(t, err)

	c, ok := out.Content().([]byte)
	assert.True(t, ok)
	assert.Equal(t, "hello from wasm", string(c))

}

func process(ctx context.Context, f *Function, m camel.Message) (camel.Message, error) {
	encoded, err := serdes.EncodeMessage(m)
	if err != nil {
		return nil, err
	}

	data, err := f.Invoke(ctx, encoded)
	if err != nil {
		return nil, err
	}

	return serdes.DecodeMessage(data)
}
