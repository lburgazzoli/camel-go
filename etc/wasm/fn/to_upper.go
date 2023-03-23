// nolint
package main

import (
	"reflect"
	"strings"
	"unsafe"

	"github.com/lburgazzoli/camel-go/pkg/wasm/interop"
	karmem "karmem.org/golang"
)

// main is required for TinyGo to compile to Wasm.
func main() {}

//export process
func _process(ptr uint32, size uint32) uint64 {
	in := ptrToMessage(ptr, size)

	in.Content = []byte(strings.ToUpper(string(in.Content)))

	p, s := messageToPtr(in)

	return (uint64(p) << uint64(32)) | uint64(s)
}

func ptrToMessage(ptr uint32, size uint32) interop.Message {
	data := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  uintptr(size), // Tinygo requires these as uintptrs even if they are int fields.
		Cap:  uintptr(size), // ^^ See https://github.com/tinygo-org/tinygo/issues/1284
	}))

	return interop.DecodeMessage(data)
}

func messageToPtr(msg interop.Message) (uint32, uint32) {
	w := karmem.NewWriter(1024)
	if _, err := msg.WriteAsRoot(w); err != nil {
		panic(err)
	}

	buf := w.Bytes()
	ptr := &buf[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))

	return uint32(unsafePtr), uint32(len(buf))
}
