package wasm

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/pkg/errors"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/wasm/serdes"
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

func (r *Runtime) Load(ctx context.Context, path string) (*Processor, error) {
	content, err := os.ReadFile(path)
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

	p := Processor{
		m:      module,
		f:      module.ExportedFunction("process"),
		malloc: module.ExportedFunction("malloc"),
		free:   module.ExportedFunction("free"),
	}

	return &p, nil
}

type Processor struct {
	m      wapi.Module
	f      wapi.Function
	malloc wapi.Function
	free   wapi.Function
}

func (p *Processor) Process(ctx context.Context, m api.Message) (api.Message, error) {
	encoded, err := serdes.Encode(m)
	if err != nil {
		return nil, err
	}

	encodedSize := uint64(len(encoded))

	results, err := p.malloc.Call(ctx, encodedSize)
	if err != nil {
		return nil, err
	}

	encodedPtr := results[0]

	defer func() { _, _ = p.free.Call(ctx, encodedPtr) }()

	if ok := p.m.Memory().Write(uint32(encodedPtr), encoded); !ok {
		return nil, fmt.Errorf(
			"memory.Write(%d, %d) out of range of memory size %d",
			encodedPtr,
			encodedSize,
			p.m.Memory().Size())
	}

	ptrSize, err := p.f.Call(ctx, encodedPtr, encodedSize)
	if err != nil {
		return nil, err
	}

	// Note: This pointer is still owned by TinyGo, so don't try to free it!
	outPtr := uint32(ptrSize[0] >> 32)
	outSize := uint32(ptrSize[0])

	bytes, ok := p.m.Memory().Read(outPtr, outSize)
	if !ok {
		return nil, fmt.Errorf(
			"memory.Read(%d, %d) out of range of memory size %d",
			outPtr,
			outSize,
			p.m.Memory().Size())
	}

	return serdes.Decode(bytes)
}
