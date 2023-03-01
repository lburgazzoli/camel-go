// nolint
package main

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/lburgazzoli/camel-go/pkg/wasm/interop"
)

// main is required for TinyGo to compile to Wasm.
func main() {}

// process message
func process(in *interop.Message) {
	fmt.Println("Processing message ", in.ID)
}

//export process
func _process(ptr uint32, size uint32) uint64 {
	in := ptrToMessage(ptr, size)

	process(&in)

	return 0
}

func ptrToMessage(ptr uint32, size uint32) interop.Message {
	data := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  uintptr(size), // Tinygo requires these as uintptrs even if they are int fields.
		Cap:  uintptr(size), // ^^ See https://github.com/tinygo-org/tinygo/issues/1284
	}))

	return interop.DecodeMessage(data)
}
