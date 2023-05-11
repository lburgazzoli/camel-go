package otherwise

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

func New(opts ...OptionFn) *Otherwise {
	answer := &Otherwise{
		DefaultStepsVerticle: processors.NewDefaultStepsVerticle(),
	}

	for _, o := range opts {
		o(answer)
	}

	return answer
}

type Otherwise struct {
	processors.DefaultStepsVerticle `yaml:",inline"`
	pid                             *actor.PID
}

func (o *Otherwise) Reify(ctx context.Context) (camel.Verticle, error) {
	c := camel.ExtractContext(ctx)
	o.DefaultVerticle.SetContext(c)

	return o, nil
}
