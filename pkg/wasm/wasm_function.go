package wasm

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/tetratelabs/wazero/api"
)

type Function struct {
	module *Module
	fn     api.Function
}

func (p *Function) invoke(ctx context.Context, in any, out any) error {
	data, err := json.Marshal(in)
	if err != nil {
		return err
	}

	// clean up the buffer
	p.module.stdin.Reset()
	p.module.stdout.Reset()

	defer func() {
		// clean up the buffer when the method
		p.module.stdin.Reset()
		p.module.stdout.Reset()
	}()

	ws, err := p.module.stdin.Write(data)
	if err != nil {
		return err
	}

	ptrSize, err := p.fn.Call(ctx, uint64(ws))
	if err != nil {
		return err
	}

	resFlag := uint32(ptrSize[0] >> 32)
	resSize := uint32(ptrSize[0])

	bytes := make([]byte, resSize)
	_, err = p.module.stdout.Read(bytes)
	if err != nil {
		return err
	}

	switch resFlag {
	case 1:
		return errors.New(string(bytes))
	default:
		return json.Unmarshal(bytes, &out)
	}
}
