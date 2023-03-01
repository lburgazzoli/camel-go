////go:build steps_transform || steps_all

package transform

import (
	"context"
	"os"
	"path"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/util/registry"
	"github.com/lburgazzoli/camel-go/pkg/wasm/serdes"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/lburgazzoli/camel-go/pkg/wasm"
)

const TAG = "transform"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &Transform{
			Identity: uuid.New(),
		}
	}
}

type Transform struct {
	camel.Identifiable
	camel.WithOutputs

	Identity string `yaml:"id"`
	Language `yaml:",inline"`

	context camel.Context
	runtime *wasm.Runtime

	processor *wasmProcessor
}

type Language struct {
	Wasm *LanguageWasm `yaml:"wasm,omitempty"`
}

type LanguageWasm struct {
	Path  string `yaml:"path"`
	Image string `yaml:"image,omitempty"`
}

func (t *Transform) ID() string {
	return t.Identity
}

func (t *Transform) Reify(ctx context.Context, camelContext camel.Context) (string, error) {

	if t.Wasm == nil {
		return "", camelerrors.MissingParameterf("wasm", "failure processing %s", TAG)
	}
	if t.Wasm.Path == "" {
		return "", camelerrors.MissingParameterf("wasm.path", "failure processing %s", TAG)
	}

	rootPath := ""

	if t.Wasm.Image != "" {
		fp, err := registry.Pull(ctx, t.Language.Wasm.Image)
		if err != nil {
			return "", err
		}

		rootPath = fp
	}

	defer func() {
		if rootPath != "" {
			_ = os.RemoveAll(rootPath)
		}
	}()

	fd, err := os.Open(path.Join(rootPath, t.Wasm.Path))
	if err != nil {
		return "", err
	}

	defer func() { _ = fd.Close() }()

	r, err := wasm.NewRuntime(ctx, wasm.Options{})
	if err != nil {
		return "", err
	}

	f, err := r.Load(ctx, "process", fd)
	if err != nil {
		return "", err
	}

	t.runtime = r
	t.context = camelContext
	t.processor = &wasmProcessor{f: f}

	return t.Identity, camelContext.Spawn(t)
}

func (t *Transform) Receive(c actor.Context) {
	msg, ok := c.Message().(camel.Message)
	if ok {
		annotations := msg.Annotations()

		out, err := t.processor.Process(context.Background(), msg)
		if err != nil {
			panic(err)
		}

		// temporary override annotations
		out.SetAnnotations(annotations)

		for _, pid := range t.Outputs() {
			if err := t.context.Send(pid, out); err != nil {
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

	return serdes.DecodeMessage(data)
}
