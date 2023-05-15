// //go:build steps_transform || steps_all

package transform

import (
	"context"

	"github.com/lburgazzoli/camel-go/pkg/core/language"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
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

	p, err := t.Language.Processor(ctx, camelContext)
	if err != nil {
		return nil, err
	}

	t.processor = p

	return t, nil
}

func (t *Transform) Receive(ac actor.Context) {
	msg, ok := ac.Message().(camel.Message)
	if ok {
		ctx := camel.Wrap(context.Background(), t.Context())

		err := t.processor(ctx, msg)
		if err != nil {
			panic(err)
		}

		ac.Request(ac.Sender(), msg)
	}
}
