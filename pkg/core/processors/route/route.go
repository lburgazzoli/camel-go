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
}

func (r *Route) Reify(ctx context.Context) (camel.Verticle, error) {
	r.SetContext(camel.GetContext(ctx))

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

		var last *actor.PID

		for s := len(items) - 1; s >= 0; s-- {
			item := items[s]

			pid, err := verticles.Spawn(ac, item)
			if err != nil {
				panic(errors.Wrapf(err, "unable to spawn verticle with id %s", item.ID()))
			}

			if last != nil {
				item.Next(last)
			}

			last = pid
		}

		if last != nil {
			r.From.Next(last)
		}

		consumer, err := r.From.Endpoint.Consumer(r.Context())
		if err != nil {
			panic(errors.Wrapf(err, "error creating consumer"))
		}

		_, err = verticles.Spawn(ac, consumer)
		if err != nil {
			panic(errors.Wrapf(err, "unable to spawn verticle with id %s", consumer.ID()))
		}

		r.Context().Registry().Set(consumer.ID(), consumer)
		r.Context().Registry().Set(r.ID(), ac.Self())
	case camel.Message:
		for _, id := range r.From.Outputs() {
			ac.Send(id, msg)
		}
	}
}
