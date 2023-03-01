package functions

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/lburgazzoli/camel-go/pkg/wasm"
	"github.com/lburgazzoli/camel-go/pkg/wasm/interop"
	"github.com/tetratelabs/wazero/api"
	karmem "karmem.org/golang"
)

func HTTPRequest(ctx context.Context, m api.Module, offset uint32, byteCount uint32) uint64 {
	buf, ok := m.Memory().Read(offset, byteCount)
	if !ok {
		panic(fmt.Errorf(
			"memory.Read(%d, %d) out of range of memory size %d",
			offset,
			byteCount,
			m.Memory().Size()))
	}

	req := interop.DecodeHttpRequest(buf)

	var httpr *http.Request

	if len(req.Params) == 0 {
		hreq, err := http.NewRequest(req.Method, req.URL, bytes.NewReader(req.Content))
		if err != nil {
			panic(err)
		}

		httpr = hreq
	} else {
		data := url.Values{}
		for i := range req.Params {
			data.Set(req.Params[i].Key, req.Params[i].Val)
		}

		hreq, err := http.NewRequest(req.Method, req.URL, strings.NewReader(data.Encode()))
		if err != nil {
			panic(err)
		}

		httpr = hreq

	}

	for i := range req.Headers {
		httpr.Header.Set(req.Headers[i].Key, req.Headers[i].Val)
	}

	client := &http.Client{}
	httpResp, err := client.Do(httpr)
	if err != nil {
		panic(err)
	}

	defer func() {
		if httpResp != nil {
			_ = httpResp.Body.Close()
		}
	}()

	res := interop.HttpResponse{}

	if httpResp != nil {
		res.Code = int32(httpResp.StatusCode)
	}

	if httpResp != nil && httpResp.Body != nil {
		all, err := io.ReadAll(httpResp.Body)
		if err != nil {
			panic(err)
		}

		res.Content = all
	}

	writer := karmem.NewWriter(1024)

	if _, err := res.WriteAsRoot(writer); err != nil {
		panic(err)
	}

	d := writer.Bytes()
	ptr, err := wasm.WriteMemory(ctx, m, d)
	if err != nil {
		panic(err)
	}

	return (ptr << uint64(32)) | uint64(len(d))
}
