package route

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/core/verticles"
	"github.com/pkg/errors"

	camel "github.com/lburgazzoli/camel-go/pkg/api"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

const TAG = "route"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &Route{
			DefaultVerticle: processors.NewDefaultVerticle(),
		}
	}
}

type Route struct {
	processors.DefaultVerticle

	Group string `yaml:"group,omitempty"`
	From  From   `yaml:"from"`

	consumerPID *actor.PID
}

func (r *Route) Reify(ctx context.Context) (camel.Verticle, error) {
	r.SetContext(camel.ExtractContext(ctx))

	return r, nil
}

func (r *Route) Receive(ac actor.Context) {
	switch msg := ac.Message().(type) {
	case *actor.Started:
		c := verticles.NewContext(r.Context(), ac)

		items, err := processors.ReifySteps(c, r.From.Steps)
		if err != nil {
			panic(errors.Wrapf(err, "error creating from steps"))
		}

		for s := range items {
			item := items[s]

			pid, err := verticles.Spawn(ac, item)
			if err != nil {
				panic(errors.Wrapf(err, "unable to spawn verticle with id %s", item.ID()))
			}

			r.Add(pid)
		}

		consumer, err := r.From.Endpoint.Consumer(r.Context())
		if err != nil {
			panic(errors.Wrapf(err, "error creating consumer"))
		}

		r.consumerPID, err = verticles.Spawn(ac, consumer)
		if err != nil {
			panic(errors.Wrapf(err, "unable to spawn verticle with id %s", consumer.ID()))
		}

		// consumer send message to the route, which route it to the
		// route steps
		consumer.Output(ac.Self())

		r.Context().Registry().Set(r.ID(), ac.Self())
	case camel.Message:
		completed := r.Dispatch(ac, msg)

		// once completed, send the message to the consumer
		if completed {
			ac.Send(r.consumerPID, msg)
		}
	}
}
