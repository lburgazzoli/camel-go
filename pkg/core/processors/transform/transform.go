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
		return New()
	}
}

func New(opts ...OptionFn) *Transform {
	answer := &Transform{
		DefaultVerticle: processors.NewDefaultVerticle(),
	}

	for _, o := range opts {
		o(answer)
	}

	return answer
}

type Transform struct {
	processors.DefaultVerticle `yaml:",inline"`

	language.Language `yaml:",inline"`

	processor camel.Processor
}

func (t *Transform) ID() string {
	return t.Identity
}

func (t *Transform) Reify(ctx context.Context) (camel.Verticle, error) {
	camelContext := camel.ExtractContext(ctx)

	t.SetContext(camelContext)

	switch {
	case t.Wasm != nil:
		p, err := t.Wasm.Processor(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		t.processor = p

	case t.Mustache != nil:
		p, err := t.Mustache.Processor(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		t.processor = p

	case t.Jq != nil:
		p, err := t.Jq.Processor(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		t.processor = p

	case t.Constant != nil:
		p, err := t.Constant.Processor(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		t.processor = p
	default:
		return nil, camelerrors.MissingParameterf("wasm || mustache || jq", "failure processing %s", TAG)

	}

	return t, nil
}

func (t *Transform) Receive(ac actor.Context) {
	msg, ok := ac.Message().(camel.Message)
	if ok {
		annotations := msg.Annotations()
		ctx := camel.Wrap(context.Background(), t.Context())

		err := t.processor(ctx, msg)
		if err != nil {
			panic(err)
		}

		// temporary override annotations
		msg.SetAnnotations(annotations)

		ac.Request(ac.Sender(), msg)
	}
}
