// nolint
package main

import (
	"reflect"
	"strings"
	"unsafe"

	karmem "karmem.org/golang"

	"github.com/lburgazzoli/camel-go/pkg/wasm/interop"
)

// main is required for TinyGo to compile to Wasm.
func main() {}

// log a message to the console using _log.
func http(request interop.HttpRequest) interop.HttpResponse {

	//_http(ptr, size)

	return interop.HttpResponse{}
}

//
//go:wasm-module camel
//export http
func _http(ptr uint32, size uint32) uint64

const template = `
{
	"channel": "{{channel}}",
	"blocks": [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "{{text}}"
			}
		}
	]
}
`

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

	content := strings.ReplaceAll(template, "{{text}}", string(in.Content))
	content = strings.ReplaceAll(template, "{{channel}}", channel)

	req := interop.NewHttpRequest()
	req.URL = "https://slack.com/api/chat.postMessage"
	req.Method = "POST"
	req.Headers = append(req.Headers, interop.Pair{Key: "Authorization", Val: "Bearer " + token})
	req.Content = []byte(content)

	_ := http(req)

}

//export process
func _process(ptr uint32, size uint32) uint64 {
	in := ptrToMessage(ptr, size)

	out := in

	process(&out)

	p, s := messageToPtr(out)

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
