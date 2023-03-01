package wasm

import (
	"context"
	"fmt"

	wapi "github.com/tetratelabs/wazero/api"
)

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
