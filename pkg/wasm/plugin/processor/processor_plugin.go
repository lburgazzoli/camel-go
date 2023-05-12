//go:build tinygo.wasm

package processor

import (
	context "context"
	wasm "github.com/knqyf263/go-plugin/wasm"
)

var processors Processors

func RegisterProcessors(p Processors) {
	processors = p
}

//export process
func _process(ptr, size uint32) uint64 {
	b := wasm.PtrToByte(ptr, size)
	req := new(Message)
	if err := req.UnmarshalVT(b); err != nil {
		return 0
	}
	response, err := processors.Process(context.Background(), req)
	if err != nil {
		ptr, size = wasm.ByteToPtr([]byte(err.Error()))
		return (uint64(ptr) << uint64(32)) | uint64(size) |
			// Indicate that this is the error string by setting the 32-th bit, assuming that
			// no data exceeds 31-bit size (2 GiB).
			(1 << 31)
	}

	b, err = response.MarshalVT()
	if err != nil {
		return 0
	}
	ptr, size = wasm.ByteToPtr(b)
	return (uint64(ptr) << uint64(32)) | uint64(size)
}
