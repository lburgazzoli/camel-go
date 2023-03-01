// nolint
package main

import (
	"reflect"
	"strconv"
	"unsafe"

	"github.com/lburgazzoli/camel-go/pkg/wasm/interop"
)

// main is required for TinyGo to compile to Wasm.
func main() {}

// log a message to the console using _log.
func http(request *interop.HttpRequest) *interop.HttpResponse {

	p, s := interop.ToPtr(request)

	ps := _http(p, s)

	outPtr := uint32(ps >> 32)
	outSize := uint32(ps)

	return ptrToResponse(outPtr, outSize)
}

//
//go:wasm-module camel
//export http
func _http(ptr uint32, size uint32) uint64

// process message
func process(in *interop.Message) {
	token := ""
	channel := ""

	for _, a := range in.Annotations {
		switch a.Key {
		case "slack.token":
			token = a.Val
		case "slack.channel":
			channel = a.Val
		}
	}

	req := interop.NewHttpRequest()
	req.URL = "https://slack.com/api/chat.postMessage"
	req.Method = "POST"
	req.Headers = append(req.Headers, interop.Pair{Key: "Authorization", Val: "Bearer " + token})
	req.Headers = append(req.Headers, interop.Pair{Key: "Content-Type", Val: "application/json"})
	req.Params = append(req.Params, interop.Pair{Key: "text", Val: string(in.Content)})
	req.Params = append(req.Params, interop.Pair{Key: "channel", Val: channel})

	res := http(&req)

	in.Content = res.Content
	in.Annotations = append(in.Annotations, interop.Pair{Key: "http.code", Val: strconv.FormatInt(int64(res.Code), 10)})

	for i := range res.Headers {
		in.Annotations = append(in.Annotations, res.Headers[i])
	}

}

//export process
func _process(ptr uint32, size uint32) uint64 {
	in := ptrToMessage(ptr, size)

	process(in)

	p, s := interop.ToPtr(in)

	return (uint64(p) << uint64(32)) | uint64(s)
}

func ptrToMessage(ptr uint32, size uint32) *interop.Message {
	data := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  uintptr(size), // Tinygo requires these as uintptrs even if they are int fields.
		Cap:  uintptr(size), // ^^ See https://github.com/tinygo-org/tinygo/issues/1284
	}))

	r := interop.DecodeMessage(data)

	return &r
}

func ptrToResponse(ptr uint32, size uint32) *interop.HttpResponse {
	data := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  uintptr(size), // Tinygo requires these as uintptrs even if they are int fields.
		Cap:  uintptr(size), // ^^ See https://github.com/tinygo-org/tinygo/issues/1284
	}))

	r := interop.DecodeHttpResponse(data)

	return &r
}
