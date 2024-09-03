package wasm

import (
	"context"
	"errors"
	"fmt"
	"os"

	wz "github.com/tetratelabs/wazero"
	wzapi "github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/sys"
)

func NewModule(ctx context.Context, r wz.Runtime, code wz.CompiledModule) (*Module, error) {
	answer := Module{}
	answer.wz = r

	config := wz.NewModuleConfig()
	config = config.WithStdout(os.Stdout)
	config = config.WithStdin(os.Stdin)
	config = config.WithStderr(os.Stderr)

	module, err := answer.wz.InstantiateModule(ctx, code, config)
	if err != nil {
		// Note: Most compilers do not exit the module after running "_start",
		// unless there was an Error. This allows you to call exported functions.
		var syse *sys.ExitError

		if ok := errors.As(err, &syse); ok && (*syse).ExitCode() != 0 {
			return nil, fmt.Errorf("unexpected exit_code: %d", (*syse).ExitCode())
		} else if !ok {
			return nil, err
		}
	}

	answer.fnAlloc = module.ExportedFunction(allocFunctionNAme)
	if answer.fnAlloc == nil {
		return nil, fmt.Errorf("mandatory function '%s' is not exported", allocFunctionNAme)
	}

	answer.fnDealloc = module.ExportedFunction(deallocFunctionNAme)
	if answer.fnAlloc == nil {
		return nil, fmt.Errorf("mandatory function '%s' is not exported", deallocFunctionNAme)
	}

	answer.module = module

	return &answer, nil
}

type Module struct {
	wz        wz.Runtime
	module    wzapi.Module
	fnAlloc   wzapi.Function
	fnDealloc wzapi.Function
}

func (m *Module) alloc(ctx context.Context, size uint64) (uint64, error) {
	results, err := m.fnAlloc.Call(ctx, size)
	if err != nil {
		return 0, fmt.Errorf("unable to allocate %d bytes, %w", size, err)
	}

	return results[0], nil
}

//nolint:unused
func (m *Module) dealloc(ctx context.Context, ptr uint64, size uint64) error {
	_, err := m.fnDealloc.Call(ctx, ptr, size)
	if err != nil {
		return fmt.Errorf("unable to deallocate %d bytes at ptr %d, %w", size, ptr, err)
	}

	return nil
}

func (m *Module) write(ctx context.Context, data []byte) (uint64, uint64, error) {
	size := uint64(len(data))

	ptr, err := m.alloc(ctx, size)
	if err != nil {
		return noPointer, noSize, err
	}

	//nolint:gosec
	if !m.module.Memory().Write(uint32(ptr), data) {
		err := fmt.Errorf(
			"memory.Write(%d, %d) out of range of memory size %d",
			ptr,
			len(data),
			m.module.Memory().Size(),
		)

		return noPointer, noSize, err
	}

	return ptr, size, nil
}

func (m *Module) Memory() wzapi.Memory {
	if m.module == nil {
		return nil
	}

	return m.module.Memory()
}

func (m *Module) Close(ctx context.Context) error {
	if m.module == nil {
		return nil
	}

	return m.module.Close(ctx)
}

func (m *Module) Processor(_ context.Context, name string) (*Processor, error) {
	if m.module == nil {
		return nil, nil
	}

	fn := m.module.ExportedFunction(name)
	if fn == nil {
		return nil, fmt.Errorf("function %s is not exported", name)
	}

	p := Processor{
		Function: Function{
			module: m,
			fn:     fn,
		},
	}

	return &p, nil
}

func (m *Module) Predicate(_ context.Context, name string) (*Predicate, error) {
	if m.module == nil {
		return nil, nil
	}

	fn := m.module.ExportedFunction(name)
	if fn == nil {
		return nil, fmt.Errorf("function %s is not exported", name)
	}

	p := Predicate{
		Function: Function{
			module: m,
			fn:     fn,
		},
	}

	return &p, nil
}
