////go:build steps_transform || steps_all

package transform

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
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
	processor languageProcessor
}

type Language struct {
	Wasm     *LanguageWasm     `yaml:"wasm,omitempty"`
	Mustache *LanguageMustache `yaml:"mustache,omitempty"`
}

func (t *Transform) ID() string {
	return t.Identity
}

func (t *Transform) Reify(ctx context.Context, camelContext camel.Context) (string, error) {

	t.context = camelContext

	switch {
	case t.Wasm != nil:
		p, err := newWasmProcessor(ctx, t.Wasm)
		if err != nil {
			return "", err
		}

		t.processor = p

	case t.Mustache != nil:
		p, err := newMustacheProcessor(ctx, t.Mustache)
		if err != nil {
			return "", err
		}

		t.processor = p
	default:
		return "", camelerrors.MissingParameterf("wasm || mustache", "failure processing %s", TAG)

	}

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

type languageProcessor interface {
	Process(context.Context, camel.Message) (camel.Message, error)
}
