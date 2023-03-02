package wasm

import (
	"context"
	"io"
	"io/fs"

	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"go.uber.org/multierr"

	"github.com/tetratelabs/wazero"
)

func NewRuntime(ctx context.Context, opt Options) (*Runtime, error) {
	cache := wazero.NewCompilationCache()

	config := wazero.NewModuleConfig()

	if opt.Stdout != nil {
		config = config.WithStdout(opt.Stdout)
	}
	if opt.Stderr != nil {
		config = config.WithStderr(opt.Stderr)
	}
	if opt.FS != nil {
		config = config.WithFS(opt.FS)
	}

	runtime := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().
		WithCompilationCache(cache))

	if _, err := wasi.NewBuilder(runtime).Instantiate(ctx); err != nil {
		return nil, err
	}

	return &Runtime{
			wz:     runtime,
			cache:  cache,
			config: config},
		nil
}

type Options struct {
	Stdout io.Writer
	Stderr io.Writer
	FS     fs.FS
}

type Runtime struct {
	wz     wazero.Runtime
	cache  wazero.CompilationCache
	config wazero.ModuleConfig
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

func (r *Runtime) Load(ctx context.Context, name string, reader io.Reader) (*Function, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	code, err := r.wz.CompileModule(ctx, content)
	if err != nil {
		return nil, err
	}

	f := Function{
		r:    r,
		cm:   code,
		name: name,
	}

	return &f, nil
}
