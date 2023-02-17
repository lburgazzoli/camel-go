package wasm

import (
	"context"

	"github.com/tetratelabs/wazero"
)

type Runtime struct {
	r wazero.Runtime
}

func NewRuntime(ctx context.Context) Runtime {
	return Runtime{
		r: wazero.NewRuntime(ctx),
	}
}
