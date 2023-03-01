package wasm

import (
	"context"
	"fmt"
	"io"
	"io/fs"

	"github.com/pkg/errors"

	wapi "github.com/tetratelabs/wazero/api"
	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	wsys "github.com/tetratelabs/wazero/sys"
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

	module, err := r.wz.InstantiateModule(ctx, code, r.config)
	if err != nil {
		// Note: Most compilers do not exit the module after running "_start",
		// unless there was an Error. This allows you to call exported functions.
		var exitErr wsys.ExitError

		if errors.As(err, &exitErr) && exitErr.ExitCode() != 0 {
			return nil, fmt.Errorf("unexpected exit_code: %d", exitErr.ExitCode())
		}

		return nil, err
	}

	p := Function{
		m:      module,
		f:      module.ExportedFunction(name),
		malloc: module.ExportedFunction("malloc"),
		free:   module.ExportedFunction("free"),
	}

	return &p, nil
}

type Function struct {
	m      wapi.Module
	f      wapi.Function
	malloc wapi.Function
	free   wapi.Function
}

func (f *Function) write(ptr uint32, data []byte) error {
	if ok := f.m.Memory().Write(ptr, data); !ok {
		return fmt.Errorf(
			"memory.Write(%d, %d) out of range of memory size %d",
			ptr,
			len(data),
			f.m.Memory().Size())
	}

	return nil
}

func (f *Function) read(ptr uint32, size uint32) ([]byte, error) {
	bytes, ok := f.m.Memory().Read(ptr, size)
	if !ok {
		return nil, fmt.Errorf(
			"memory.Read(%d, %d) out of range of memory size %d",
			ptr,
			size,
			f.m.Memory().Size())
	}

	return bytes, nil
}

func (f *Function) Invoke(ctx context.Context, data []byte) ([]byte, error) {

	dataSize := uint64(len(data))

	results, err := f.malloc.Call(ctx, dataSize)
	if err != nil {
		return nil, err
	}

	dataPtr := results[0]

	defer func() { _, _ = f.free.Call(ctx, dataPtr) }()

	if err := f.write(uint32(dataPtr), data); err != nil {
		return nil, err
	}

	ptrSize, err := f.f.Call(ctx, dataPtr, dataSize)
	if err != nil {
		return nil, err
	}

	if ptrSize[0] == 0 {
		return nil, nil
	}

	// Note: This pointer is still owned by TinyGo, so don't try to free it!
	outPtr := uint32(ptrSize[0] >> 32)
	outSize := uint32(ptrSize[0])

	return f.read(outPtr, outSize)
}
