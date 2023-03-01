// nolint
package main

import (
	"reflect"
	"unsafe"

	"github.com/lburgazzoli/camel-go/pkg/wasm/interop"
	karmem "karmem.org/golang"
)

// main is required for TinyGo to compile to Wasm.
func main() {}

// process message
func process(in *interop.Message) {

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

	reader := karmem.NewReader(data)
	decoded := interop.NewMessageViewer(reader, 0)

	out := interop.Message{
		ID:            decoded.ID(reader),
		Source:        decoded.Source(reader),
		Type:          decoded.Type(reader),
		Subject:       decoded.Subject(reader),
		ContentType:   decoded.ContentType(reader),
		ContentSchema: decoded.ContentSchema(reader),
		Time:          decoded.Time(),
		Content:       decoded.Content(reader),
	}

	return out
}
