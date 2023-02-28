////go:build steps_transform || steps_all

package transform

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"

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

	context   camel.Context
	processor *wasm.Processor
	runtime   *wasm.Runtime
}

type Language struct {
	Wasm *LanguageWasm `yaml:"wasm,omitempty"`
}

type LanguageWasm struct {
	Path string `yaml:"path"`
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

	r, err := wasm.NewRuntime(ctx, wasm.Options{})
	if err != nil {
		return "", err
	}

	m, err := r.Load(ctx, t.Wasm.Path)
	if err != nil {
		return "", err
	}

	t.runtime = r
	t.context = camelContext
	t.processor = m

	return t.Identity, camelContext.Spawn(t)
}

func (t *Transform) Receive(c actor.Context) {
	msg, ok := c.Message().(camel.Message)
	if ok {
		out, err := t.processor.Process(context.Background(), msg)
		if err != nil {
			panic(err)
		}

		for _, pid := range t.Outputs() {
			if err := t.context.Send(pid, out); err != nil {
				panic(err)
			}
		}
	}
}
