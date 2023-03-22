package choice

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/core/verticles"
	"github.com/pkg/errors"

	"github.com/lburgazzoli/camel-go/pkg/core/language"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

type When struct {
	processors.DefaultVerticle `yaml:",inline"`
	language.Language          `yaml:",inline"`

	predicate camel.Predicate
	pid       *actor.PID

	Steps []processors.Step `yaml:"steps,omitempty"`
}

func (w *When) Reify(ctx context.Context) (camel.Verticle, error) {
	c := camel.GetContext(ctx)

	w.DefaultVerticle.SetContext(c)

	switch {
	case w.Jq != nil:
		p, err := w.Jq.Predicate(ctx, c)
		if err != nil {
			return nil, err
		}

		w.predicate = p
	default:
		return nil, camelerrors.MissingParameterf("jq", "failure processing %s", TAG)
	}

	return w, nil
}

func (w *When) Matches(ctx context.Context, msg camel.Message) (bool, error) {
	if w.predicate == nil {
		return false, camelerrors.InternalErrorf("not configured")
	}

	return w.predicate(ctx, msg)
}

func (w *When) Receive(ac actor.Context) {
	switch msg := ac.Message().(type) {
	case *actor.Started:
		ctx := verticles.NewContext(w.Context(), ac)

		items, err := processors.ReifySteps(ctx, w.Steps)
		if err != nil {
			panic(errors.Wrapf(err, "error creating when steps"))
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
			w.Next(last)
		}
	case camel.Message:
		w.Dispatch(ac, msg)
	}
}
