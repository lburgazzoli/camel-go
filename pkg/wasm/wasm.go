package wasm

import (
	"context"
	"io"

	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"go.uber.org/multierr"

	"github.com/tetratelabs/wazero"
)

func NewRuntime(ctx context.Context) (*Runtime, error) {
	cache := wazero.NewCompilationCache()

	runtime := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().
		WithCompilationCache(cache))

	if _, err := wasi.NewBuilder(runtime).Instantiate(ctx); err != nil {
		return nil, err
	}

	return &Runtime{
			wz:    runtime,
			cache: cache},
		nil
}

type Runtime struct {
	wz    wazero.Runtime
	cache wazero.CompilationCache
}

func (r *Runtime) Export(ctx context.Context, name string, fn interface{}) error {
	_, err := r.wz.NewHostModuleBuilder("camel").
		NewFunctionBuilder().
		WithFunc(fn).
		Export(name).
		Instantiate(ctx)

	return err
}

func (r *Runtime) Close(ctx context.Context) error {
	var err error

	if r.wz != nil {
		if e := r.wz.Close(ctx); e != nil {
			err = multierr.Append(err, e)
		}
	}
	if r.cache != nil {
		if e := r.cache.Close(ctx); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}

func (r *Runtime) Load(ctx context.Context, in io.ReadCloser) (*Module, error) {
	content, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}

	// Compile the WebAssembly module using the default configuration.
	code, err := r.wz.CompileModule(ctx, content)
	if err != nil {
		return nil, err
	}

	return NewModule(ctx, r.wz, code)
}
