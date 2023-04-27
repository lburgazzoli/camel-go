package choice

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/core/verticles"
	"github.com/pkg/errors"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

type Otherwise struct {
	processors.DefaultVerticle `yaml:",inline"`

	pid *actor.PID

	Steps []processors.Step `yaml:"steps,omitempty"`
}

func (o *Otherwise) Reify(ctx context.Context) (camel.Verticle, error) {
	c := camel.ExtractContext(ctx)
	o.DefaultVerticle.SetContext(c)

	return o, nil
}

func (o *Otherwise) Receive(ac actor.Context) {
	switch msg := ac.Message().(type) {
	case *actor.Started:
		ctx := verticles.NewContext(o.Context(), ac)

		items, err := processors.ReifySteps(ctx, o.Steps)
		if err != nil {
			panic(errors.Wrapf(err, "error creating when steps"))
		}

		for s := range items {
			item := items[s]

			pid, err := verticles.Spawn(ac, item)
			if err != nil {
				panic(errors.Wrapf(err, "unable to spawn verticle with id %s", item.ID()))
			}

			o.Add(pid)
		}
	case camel.Message:
		completed := o.Dispatch(ac, msg)

		// once completed, send the message to the sender
		if completed {
			ac.Send(ac.Sender(), branchDone{M: msg})
		}
	}
}
