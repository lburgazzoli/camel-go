// //go:build steps_transform || steps_all

package transform

import (
	"context"

	"github.com/lburgazzoli/camel-go/pkg/core/language"

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

	language.Language `yaml:",inline"`

	processor camel.Processor
}

func (t *Transform) ID() string {
	return t.Identity
}

func (t *Transform) Reify(ctx context.Context, camelContext camel.Context) (string, error) {

	t.SetContext(camelContext)

	switch {
	case t.Wasm != nil:
		p, err := t.Wasm.Processor(ctx, camelContext)
		if err != nil {
			return "", err
		}

		t.processor = p

	case t.Mustache != nil:
		p, err := t.Mustache.Processor(ctx, camelContext)
		if err != nil {
			return "", err
		}

		t.processor = p

	case t.Jq != nil:
		p, err := t.Jq.Processor(ctx, camelContext)
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

		err := t.processor(context.Background(), msg)
		if err != nil {
			panic(err)
		}

		// temporary override annotations
		msg.SetAnnotations(annotations)

		t.Dispatch(msg)
	}
}
