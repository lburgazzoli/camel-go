// //go:build steps_process || steps_all

package setbody

import (
	"context"

	"github.com/lburgazzoli/camel-go/pkg/core/processors/transform"

	"github.com/lburgazzoli/camel-go/pkg/core/language"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

const TAG = "setBody"

func init() {
	processors.Types[TAG] = func() interface{} {
		return transform.New()
	}
}

func New(opts ...OptionFn) *SetBody {
	answer := &SetBody{
		DefaultVerticle: processors.NewDefaultVerticle(),
	}

	for _, o := range opts {
		o(answer)
	}

	return answer
}

type SetBody struct {
	processors.DefaultVerticle `yaml:",inline"`
	language.Language          `yaml:",inline"`

	processor camel.Processor
}

func (s *SetBody) Reify(ctx context.Context) (camel.Verticle, error) {
	camelContext := camel.ExtractContext(ctx)

	s.SetContext(camelContext)

	p, err := s.Language.Processor(ctx, camelContext)
	if err != nil {
		return nil, err
	}

	s.processor = p

	return s, nil
}

func (s *SetBody) Receive(ac actor.Context) {
	msg, ok := ac.Message().(camel.Message)
	if ok {
		ctx := camel.Wrap(context.Background(), s.Context())

		err := s.processor(ctx, msg)
		if err != nil {
			panic(err)
		}

		ac.Request(ac.Sender(), msg)
	}
}
