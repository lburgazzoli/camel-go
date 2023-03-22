// //go:build components_wasm || components_all
package wasm

import (
	"context"
	"os"
	"path"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/util/registry"
	"github.com/lburgazzoli/camel-go/pkg/wasm"
	"github.com/lburgazzoli/camel-go/pkg/wasm/functions"
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

	if err := r.Export(ctx, "http", functions.HTTPRequest); err != nil {
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
			if err := p.context.SendTo(pid, out); err != nil {
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
