package choice

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

func NewOtherwise(steps ...processors.Step) *Otherwise {
	w := Otherwise{
		DefaultStepsVerticle: processors.NewDefaultStepsVerticle(),
	}

	w.Steps = steps

	return &w
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
