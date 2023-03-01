package wasm

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/lburgazzoli/camel-go/pkg/wasm/interop"
	"github.com/tetratelabs/wazero/api"
	karmem "karmem.org/golang"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/wasm/serdes"

	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/stretchr/testify/assert"
)

func TestWASM(t *testing.T) {
	ctx := context.Background()

	r, err := NewRuntime(ctx, Options{})
	assert.Nil(t, err)

	defer func() { _ = r.Close(ctx) }()

	fd, err := os.Open("../../etc/fn/simple_process.wasm")
	assert.Nil(t, err)

	f, err := r.Load(ctx, "process", fd)
	assert.Nil(t, err)

	in, err := message.New()
	assert.Nil(t, err)

	out, err := process(ctx, f, in)
	assert.Nil(t, err)

	c, ok := out.Content().([]byte)
	assert.True(t, ok)
	assert.Equal(t, "hello from wasm", string(c))

}

func process(ctx context.Context, f *Function, m camel.Message) (camel.Message, error) {
	encoded, err := serdes.EncodeMessage(m)
	if err != nil {
		return nil, err
	}

	data, err := f.Invoke(ctx, encoded)
	if err != nil {
		return nil, err
	}

	return serdes.DecodeMessage(data)
}

func TestCallbackWASM(t *testing.T) {
	ctx := context.Background()

	r, err := NewRuntime(ctx, Options{})
	assert.Nil(t, err)

	r.Export(ctx, "http", http_call)

	defer func() { _ = r.Close(ctx) }()

	fd, err := os.Open("../../etc/components/slack.wasm")
	assert.Nil(t, err)

	f, err := r.Load(ctx, "process", fd)
	assert.Nil(t, err)

	in, err := message.New()
	assert.Nil(t, err)

	in.SetAnnotation("slack.token", uuid.New())
	in.SetAnnotation("slack.channel", uuid.New())
	in.SetContent("hello from gamel")

	out, err := process(ctx, f, in)
	assert.Nil(t, err)

	c, ok := out.Content().([]byte)
	assert.True(t, ok)
	assert.Contains(t, string(c), "invalid_auth")

}

func http_call(ctx context.Context, m api.Module, offset uint32, byteCount uint32) uint64 {
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

	defer httpResp.Body.Close()

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
	p, err := WriteMemory(ctx, m, d)
	if err != nil {
		panic(err)
	}

	return (p << uint64(32)) | uint64(len(d))
}
