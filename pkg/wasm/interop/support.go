package interop

import (
	"unsafe"

	karmem "karmem.org/golang"
)

type Writable interface {
	WriteAsRoot(*karmem.Writer) (uint, error)
}

func DecodeMessage(data []byte) Message {
	reader := karmem.NewReader(data)
	decoded := NewMessageViewer(reader, 0)

	m := Message{
		ID:            decoded.ID(reader),
		Source:        decoded.Source(reader),
		Type:          decoded.Type(reader),
		Subject:       decoded.Subject(reader),
		ContentType:   decoded.ContentType(reader),
		ContentSchema: decoded.ContentSchema(reader),
		Time:          decoded.Time(),
		Content:       decoded.Content(reader),
	}

	decoded.Annotations(reader)

	for _, a := range decoded.Annotations(reader) {
		m.Annotations = append(m.Annotations, Pair{
			Key: a.Key(reader),
			Val: a.Val(reader),
		})
	}

	return m
}

func DecodeHttpRequest(data []byte) HttpRequest {
	reader := karmem.NewReader(data)
	decoded := NewHttpRequestViewer(reader, 0)

	r := HttpRequest{
		URL:     decoded.URL(reader),
		Method:  decoded.Method(reader),
		Content: decoded.Content(reader),
	}

	for _, h := range decoded.Headers(reader) {
		r.Headers = append(r.Headers, Pair{
			Key: h.Key(reader),
			Val: h.Val(reader),
		})
	}

	for _, p := range decoded.Params(reader) {
		r.Params = append(r.Params, Pair{
			Key: p.Key(reader),
			Val: p.Val(reader),
		})
	}

	return r
}

func DecodeHttpResponse(data []byte) HttpResponse {
	reader := karmem.NewReader(data)
	decoded := NewHttpResponseViewer(reader, 0)

	r := HttpResponse{
		Code:    decoded.Code(),
		Content: decoded.Content(reader),
	}

	for _, h := range decoded.Headers(reader) {
		r.Headers = append(r.Headers, Pair{
			Key: h.Key(reader),
			Val: h.Val(reader),
		})
	}

	return r
}

func ToPtr(writable Writable) (uint32, uint32) {
	w := karmem.NewWriter(1024)
	if _, err := writable.WriteAsRoot(w); err != nil {
		panic(err)
	}

	buf := w.Bytes()
	ptr := &buf[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))

	return uint32(unsafePtr), uint32(len(buf))
}
