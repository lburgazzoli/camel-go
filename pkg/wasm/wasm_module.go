package wasm

import (
	"context"
	"errors"

	wz "github.com/tetratelabs/wazero"
	wzapi "github.com/tetratelabs/wazero/api"
)

type Module struct {
	wz     wz.Runtime
	module wzapi.Module
}

func (m *Module) Close(ctx context.Context) error {
	if m.module == nil {
		return nil
	}

	return m.module.Close(ctx)
}

func (m *Module) Processor(ctx context.Context) (*Processor, error) {
	if m.module == nil {
		return nil, nil
	}

	malloc := m.module.ExportedFunction("malloc")
	if malloc == nil {
		return nil, errors.New("malloc is not exported")
	}

	free := m.module.ExportedFunction("free")
	if free == nil {
		return nil, errors.New("free is not exported")
	}

	fn := m.module.ExportedFunction("process")
	if fn == nil {
		return nil, errors.New("process is not exported")
	}

	p := Processor{
		Function: Function{
			module: m.module,
			malloc: malloc,
			free:   free,
			fn:     fn,
		},
	}

	return &p, nil
}
