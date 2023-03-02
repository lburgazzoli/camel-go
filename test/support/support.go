package support

import (
	"context"
	"testing"

	"github.com/lburgazzoli/camel-go/pkg/wasm"
	"github.com/lburgazzoli/camel-go/pkg/wasm/serdes"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/stretchr/testify/assert"
)

func Reify(t *testing.T, c camel.Context, r processors.Reifyable) (string, error) {
	t.Helper()

	id, err := r.Reify(context.Background(), c)
	assert.Nil(t, err)
	assert.NotNil(t, id)

	return id, err
}

func Process(ctx context.Context, f *wasm.Function, m camel.Message) (camel.Message, error) {
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
