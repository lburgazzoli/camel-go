package support

import (
	"context"

	"github.com/lburgazzoli/camel-go/pkg/wasm"
	"github.com/lburgazzoli/camel-go/pkg/wasm/serdes"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
)

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
