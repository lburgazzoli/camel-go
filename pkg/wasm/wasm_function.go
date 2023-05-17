package wasm

import (
	"context"
	"errors"

	"github.com/tetratelabs/wazero/api"
)

type VTProtoSerde interface {
	MarshalVT() ([]byte, error)
	UnmarshalVT([]byte) error
}

type Function struct {
	module *Module
	fn     api.Function
}

func (p *Function) invoke(ctx context.Context, inout VTProtoSerde) error {
	data, err := inout.MarshalVT()
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
		return inout.UnmarshalVT(bytes)
	}
}
