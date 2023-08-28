package wasm

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	wz "github.com/tetratelabs/wazero"
	wzapi "github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/sys"
)

func NewModule(ctx context.Context, r wz.Runtime, code wz.CompiledModule) (*Module, error) {
	answer := Module{}
	answer.wz = r

	config := wz.NewModuleConfig()
	config = config.WithStdout(io.Writer(&answer.stdout))
	config = config.WithStdin(io.Reader(&answer.stdin))
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

	answer.module = module

	return &answer, nil
}

type Module struct {
	wz     wz.Runtime
	module wzapi.Module
	stdin  bytes.Buffer
	stdout bytes.Buffer
}

func (m *Module) Close(ctx context.Context) error {
	if m.module == nil {
		return nil
	}

	return m.module.Close(ctx)
}

func (m *Module) Processor(_ context.Context) (*Processor, error) {
	if m.module == nil {
		return nil, nil
	}

	fn := m.module.ExportedFunction("process")
	if fn == nil {
		return nil, errors.New("process is not exported")
	}

	p := Processor{
		Function: Function{
			module: m,
			fn:     fn,
		},
	}

	return &p, nil
}

func (m *Module) Predicate(_ context.Context) (*Predicate, error) {
	if m.module == nil {
		return nil, nil
	}

	fn := m.module.ExportedFunction("test")
	if fn == nil {
		return nil, errors.New("test is not exported")
	}

	p := Predicate{
		Function: Function{
			module: m,
			fn:     fn,
		},
	}

	return &p, nil
}
