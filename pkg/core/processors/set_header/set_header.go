// //go:build steps_process || steps_all

package setheader

import (
	"context"

	"github.com/lburgazzoli/camel-go/pkg/core/language"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

const TAG = "setHeader"

func init() {
	processors.Types[TAG] = func() interface{} {
		return New()
	}
}

func New(opts ...OptionFn) *SetHeader {
	answer := &SetHeader{
		DefaultVerticle: processors.NewDefaultVerticle(),
	}

	for _, o := range opts {
		o(answer)
	}

	return answer
}

type SetHeader struct {
	processors.DefaultVerticle `yaml:",inline"`
	language.Language          `yaml:",inline"`

	Name string `yaml:"name"`

	transformer camel.Transformer
}

func (p *SetHeader) ID() string {
	return p.Identity
}

func (p *SetHeader) Reify(ctx context.Context) (camel.Verticle, error) {
	camelContext := camel.ExtractContext(ctx)

	p.SetContext(camelContext)

	t, err := p.Language.Transformer(ctx, camelContext)
	if err != nil {
		return nil, err
	}

	p.transformer = t

	return p, nil
}

func (p *SetHeader) Receive(ac actor.Context) {
	msg, ok := ac.Message().(camel.Message)
	if ok {
		ctx := camel.Wrap(context.Background(), p.Context())

		answer, err := p.transformer(ctx, msg)
		if err != nil {
			panic(err)
		}

		msg.SetHeader(p.Name, answer)

		ac.Request(ac.Sender(), msg)
	}
}
