// //go:build steps_process || steps_all

package setbody

import (
	"context"

	"github.com/lburgazzoli/camel-go/pkg/core/language"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

const TAG = "setBody"

func init() {
	processors.Types[TAG] = func() interface{} {
		return New()
	}
}

func New() *SetBody {
	return &SetBody{
		DefaultVerticle: processors.NewDefaultVerticle(),
	}
}

type SetBody struct {
	processors.DefaultVerticle `yaml:",inline"`
	language.Language          `yaml:",inline"`

	processor camel.Processor
}

func (p *SetBody) Reify(ctx context.Context) (camel.Verticle, error) {
	camelContext := camel.ExtractContext(ctx)

	p.SetContext(camelContext)

	switch {
	case p.Wasm != nil:
		proc, err := p.Wasm.Processor(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		p.processor = proc

	case p.Mustache != nil:
		proc, err := p.Mustache.Processor(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		p.processor = proc

	case p.Jq != nil:
		proc, err := p.Jq.Processor(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		p.processor = proc

	case p.Constant != nil:
		proc, err := p.Constant.Processor(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		p.processor = proc
	default:
		return nil, camelerrors.MissingParameterf("wasm || mustache || jq", "failure processing %s", TAG)

	}

	return p, nil
}

func (p *SetBody) Receive(ac actor.Context) {
	msg, ok := ac.Message().(camel.Message)
	if ok {
		annotations := msg.Annotations()
		ctx := camel.Wrap(context.Background(), p.Context())

		err := p.processor(ctx, msg)
		if err != nil {
			panic(err)
		}

		// temporary override annotations
		msg.SetAnnotations(annotations)

		ac.Request(ac.Sender(), msg)
	}
}
