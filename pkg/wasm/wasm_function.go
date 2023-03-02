package wasm

import (
	"context"
	"fmt"

	wapi "github.com/tetratelabs/wazero/api"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/sys"
)

type Function struct {
	r    *Runtime
	cm   wazero.CompiledModule
	name string
}

// Invoke invoke a function.
// TODO: this is likely be highly inefficient as the model is
// TODO: loaded for each invocation but at this stage I don't
// TODO: have enough knowledge of memory management to make
// TODO: things better.
func (f *Function) Invoke(ctx context.Context, data []byte) ([]byte, error) {

	var err error
	var dataPtr uint64
	var module wapi.Module
	var fn wapi.Function
	var free wapi.Function

	defer func() {
		if dataPtr != 0 && free != nil {
			_, _ = free.Call(ctx, dataPtr)
		}

		if module != nil {
			_ = module.Close(ctx)
		}
	}()

	module, err = f.r.wz.InstantiateModule(ctx, f.cm, f.r.config)
	if err != nil {
		// Note: Most compilers do not exit the module after running "_start",
		// unless there was an Error. This allows you to call exported functions.

		//nolint:errorlint
		if exitErr, ok := err.(*sys.ExitError); ok && exitErr.ExitCode() != 0 {
			return nil, fmt.Errorf("unexpected exit_code: %d", exitErr.ExitCode())
		}

		return nil, err
	}

	fn = module.ExportedFunction(f.name)
	if fn == nil {
		return nil, fmt.Errorf("unable to load function: %s", f.name)
	}
	free = module.ExportedFunction("free")
	if free == nil {
		return nil, fmt.Errorf("unable to load function: %s", "free")
	}

	dataPtr, err = WriteMemory(ctx, module, data)
	if err != nil {
		return nil, err
	}

	// Note: This pointer is owned by the wasm runtime, so don't try to free it!
	ptrAndSize, err := fn.Call(ctx, dataPtr, uint64(len(data)))
	if err != nil {
		return nil, err
	}

	return ReadMemory(
		ctx,
		module.Memory(),
		uint32(ptrAndSize[0]>>32),
		uint32(ptrAndSize[0]))
}
