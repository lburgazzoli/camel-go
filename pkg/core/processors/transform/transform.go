// //go:build steps_transform || steps_all

package transform

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

const TAG = "transform"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &Transform{
			DefaultVerticle: processors.NewDefaultVerticle(),
		}
	}
}

type Transform struct {
	processors.DefaultVerticle `yaml:",inline"`

	Language `yaml:",inline"`

	processor languageProcessor
}

type Language struct {
	Wasm     *LanguageWasm     `yaml:"wasm,omitempty"`
	Mustache *LanguageMustache `yaml:"mustache,omitempty"`
	Jq       *LanguageJq       `yaml:"jq,omitempty"`
}

func (t *Transform) ID() string {
	return t.Identity
}

func (t *Transform) Reify(ctx context.Context, camelContext camel.Context) (string, error) {

	t.SetContext(camelContext)

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

	case t.Jq != nil:
		p, err := newJqProcessor(ctx, t.Jq)
		if err != nil {
			return "", err
		}

		t.processor = p
	default:
		return "", camelerrors.MissingParameterf("wasm || mustache || jq", "failure processing %s", TAG)

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

		t.Dispatch(out)
	}
}

type languageProcessor interface {
	Process(context.Context, camel.Message) (camel.Message, error)
}
