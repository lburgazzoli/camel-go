package wasm

import (
	"context"
	"errors"
	"fmt"

	"github.com/tetratelabs/wazero/api"
)

type VTProtoSerde interface {
	MarshalVT() ([]byte, error)
	UnmarshalVT([]byte) error
}

type Function struct {
	module api.Module
	malloc api.Function
	free   api.Function
	fn     api.Function
}

func (p *Function) invoke(ctx context.Context, inout VTProtoSerde) error {
	data, err := inout.MarshalVT()
	if err != nil {
		return err
	}

	dataSize := uint64(len(data))

	var dataPtr uint64
	// If the input data is not empty, we must allocate the in-Wasm memory to store it, and pass to the plugin.
	if dataSize != 0 {
		results, err := p.malloc.Call(ctx, dataSize)
		if err != nil {
			return err
		}
		dataPtr = results[0]
		// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
		// So, we have to free it when finished
		defer func() {
			_, _ = p.free.Call(ctx, dataPtr)
		}()

		// The pointer is a linear memory offset, which is where we write the name.
		if !p.module.Memory().Write(uint32(dataPtr), data) {
			return fmt.Errorf("Memory.Write(%d, %d) out of range of memory size %d", dataPtr, dataSize, p.module.Memory().Size())
		}
	}

	ptrSize, err := p.fn.Call(ctx, dataPtr, dataSize)
	if err != nil {
		return err
	}

	resPtr := uint32(ptrSize[0] >> 32)
	resSize := uint32(ptrSize[0])

	var isErrResponse bool

	if (resSize & (1 << 31)) > 0 {
		isErrResponse = true
		resSize &^= 1 << 31
	}

	// We don't need the memory after deserialization: make sure it is freed.
	if resPtr != 0 {
		defer func() {
			_, _ = p.free.Call(ctx, uint64(resPtr))
		}()
	}

	// The pointer is a linear memory offset, which is where we write the name.
	bytes, ok := p.module.Memory().Read(resPtr, resSize)
	if !ok {
		return fmt.Errorf("Memory.Read(%d, %d) out of range of memory size %d",
			resPtr,
			resSize,
			p.module.Memory().Size())
	}

	if isErrResponse {
		return errors.New(string(bytes))
	}

	if err = inout.UnmarshalVT(bytes); err != nil {
		return err
	}

	return nil
}
