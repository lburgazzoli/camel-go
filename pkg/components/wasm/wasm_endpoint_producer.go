// //go:build components_wasm || components_all
package wasm

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/lburgazzoli/camel-go/pkg/wasm/interop"
	karmem "karmem.org/golang"

	"github.com/tetratelabs/wazero/api"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/util/registry"
	"github.com/lburgazzoli/camel-go/pkg/wasm"
	"github.com/lburgazzoli/camel-go/pkg/wasm/serdes"
)

type Producer struct {
	camel.WithOutputs

	id       string
	endpoint *Endpoint

	context camel.Context
	runtime *wasm.Runtime

	processor *wasmProcessor
}

func (p *Producer) ID() string {
	return p.id
}

func (p *Producer) Endpoint() camel.Endpoint {
	return p.endpoint
}

func (p *Producer) Start(ctx context.Context) error {
	rootPath := ""

	if p.endpoint.config.Image != "" {
		fp, err := registry.Pull(ctx, p.endpoint.config.Image)
		if err != nil {
			return err
		}

		rootPath = fp
	}

	defer func() {
		if rootPath != "" {
			_ = os.RemoveAll(rootPath)
		}
	}()

	fd, err := os.Open(path.Join(rootPath, p.endpoint.config.Remaining))
	if err != nil {
		return err
	}

	defer func() { _ = fd.Close() }()

	r, err := wasm.NewRuntime(ctx, wasm.Options{Stdout: os.Stdout, Stderr: os.Stderr})
	if err != nil {
		return err
	}

	if err := r.Export(ctx, "http", p.callHTTP); err != nil {
		return err
	}

	f, err := r.Load(ctx, "process", fd)
	if err != nil {
		return err
	}

	p.context = p.endpoint.Component().Context()
	p.runtime = r
	p.processor = &wasmProcessor{f: f}

	return nil
}

func (p *Producer) Stop(ctx context.Context) error {
	if p.runtime != nil {
		return p.runtime.Close(ctx)
	}

	return nil
}

func (p *Producer) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		_ = p.Start(context.Background())
	case *actor.Stopping:
		_ = p.Stop(context.Background())
	case camel.Message:
		annotations := msg.Annotations()

		// TODO: find a type safe way to propagate parameters to the engine
		//       this is just for POC
		for k, v := range p.endpoint.config.Other {
			msg.SetAnnotation(k, v)
		}

		out, err := p.processor.Process(context.Background(), msg)
		if err != nil {
			panic(err)
		}

		// temporary override annotations
		out.SetAnnotations(annotations)

		for _, pid := range p.Outputs() {
			if err := p.context.Send(pid, out); err != nil {
				panic(err)
			}
		}
	}
}

type wasmProcessor struct {
	f *wasm.Function
}

func (p *wasmProcessor) Process(ctx context.Context, m camel.Message) (camel.Message, error) {

	encoded, err := serdes.EncodeMessage(m)
	if err != nil {
		return nil, err
	}

	data, err := p.f.Invoke(ctx, encoded)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return m, nil
	}

	return serdes.DecodeMessage(data)
}

func (p *Producer) callHTTP(ctx context.Context, m api.Module, offset uint32, byteCount uint32) uint64 {
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
