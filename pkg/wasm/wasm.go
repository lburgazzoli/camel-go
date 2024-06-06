package wasm

import (
	"context"
	"io"

	"github.com/lburgazzoli/camel-go/pkg/api"
	wzapi "github.com/tetratelabs/wazero/api"
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

	r := Runtime{
		wz:    runtime,
		cache: cache,
	}

	builder := r.wz.NewHostModuleBuilder("env")

	for _, fn := range BuiltInFunctions() {
		_ = builder.NewFunctionBuilder().
			WithGoModuleFunction(
				wzapi.GoModuleFunc(func(ctx context.Context, m wzapi.Module, stack []uint64) {
					//nolint:forcetypeassert
					mod := ctx.Value(contextKeyModule).(*Module)

					//nolint:forcetypeassert
					msg := ctx.Value(contextKeyMessage).(api.Message)

					if err := fn.Fn(ctx, mod, msg, stack); err != nil {
						panic(err)
					}
				}),
				fn.Params,
				fn.Results,
			).
			Export(fn.Name)
	}

	_, err := builder.Instantiate(ctx)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

type HostFunction struct {
	Name    string
	Fn      func(ctx context.Context, mod *Module, msg api.Message, stack []uint64) error
	Params  []wzapi.ValueType
	Results []wzapi.ValueType
}

type Runtime struct {
	wz    wazero.Runtime
	cache wazero.CompilationCache
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
