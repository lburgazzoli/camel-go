package processors

import (
	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/verticles"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/pkg/errors"
)

func NewDefaultVerticle() DefaultVerticle {
	return DefaultVerticle{
		Identity: uuid.New(),
	}
}

type DefaultVerticle struct {
	camel.Identifiable

	Identity string `yaml:"id"`

	context camel.Context
}

func (v *DefaultVerticle) Context() camel.Context {
	return v.context
}

func (v *DefaultVerticle) SetContext(ctx camel.Context) {
	v.context = ctx
}

func (v *DefaultVerticle) ID() string {
	return v.Identity
}

// Dispatch send messages to the child steps, returns true if the dispatch is completed.
func (v *DefaultVerticle) Dispatch(ac actor.Context, msg camel.Message) bool {

	pids := ac.Children()

	// no children
	if len(pids) == 0 {
		return true
	}

	if !verticles.Contains(pids, ac.Sender()) {
		// this is not a message coming from a children, send to the first one
		ac.Request(pids[0], msg)

		return false
	}

	for i := range len(pids) {
		pid := pids[i]

		if pid.Equal(ac.Sender()) && i != len(pids)-1 {
			// send the message to the next one
			ac.Request(pids[i+1], msg)
			return false
		}
	}

	return true
}

// StepsDone is a marker type to help break out.
type StepsDone struct {
	M camel.Message
}

func NewDefaultStepsVerticle() DefaultStepsVerticle {
	return DefaultStepsVerticle{
		DefaultVerticle: NewDefaultVerticle(),
	}
}

type DefaultStepsVerticle struct {
	DefaultVerticle `yaml:",inline"`

	Steps []Step `yaml:"steps,omitempty"`
}

func (v *DefaultStepsVerticle) Receive(ac actor.Context) {
	switch msg := ac.Message().(type) {
	case *actor.Started:
		ctx := verticles.NewContext(v.Context(), ac)

		items, err := ReifySteps(ctx, v.Steps)
		if err != nil {
			panic(errors.Wrapf(err, "error creating step"))
		}

		for s := range items {
			item := items[s]

			_, err := verticles.Spawn(ac, item)
			if err != nil {
				panic(errors.Wrapf(err, "unable to spawn verticle with id %s", item.ID()))
			}
		}
	case camel.Message:
		completed := v.Dispatch(ac, msg)

		// once completed, send the message to the parent
		if completed {
			ac.Request(ac.Parent(), StepsDone{M: msg})
		}
	}
}
