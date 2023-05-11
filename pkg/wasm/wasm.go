package wasm

import (
	"context"
	"io"
	"io/fs"
	"os"

	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"go.uber.org/multierr"

	"github.com/tetratelabs/wazero"

	pp "github.com/lburgazzoli/camel-go/pkg/wasm/plugin/processor"
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

func (r *Runtime) Load(ctx context.Context, in io.ReadCloser) (*Function, error) {
	// Initialize a plugin loader
	p, err := pp.NewProcessorsPlugin(ctx)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}

	// Load a plugin
	plugin, err := p.Load(ctx, content)
	if err != nil {
		return nil, err
	}

	f := Function{
		processor: plugin,
	}

	return &f, nil
}

func (r *Runtime) LoadFromPath(ctx context.Context, path string) (*Function, error) {
	// Initialize a plugin loader
	p, err := pp.NewProcessorsPlugin(ctx)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Load a plugin
	plugin, err := p.Load(ctx, content)
	if err != nil {
		return nil, err
	}

	f := Function{
		processor: plugin,
	}

	return &f, nil
}
