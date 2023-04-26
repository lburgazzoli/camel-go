package processors

import (
	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

func NewDefaultVerticle() DefaultVerticle {
	return DefaultVerticle{
		Identity: uuid.New(),
		pids:     actor.NewPIDSet(),
	}
}

type DefaultVerticle struct {
	camel.Identifiable
	camel.WithOutputs

	Identity string `yaml:"id"`

	context camel.Context
	pids    *actor.PIDSet
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

func (v *DefaultVerticle) Add(pid *actor.PID) {
	if pid == nil {
		return
	}

	v.pids.Add(pid)
}

// Dispatch send messages to the child steps, returns true if the dispatch is completed
func (v *DefaultVerticle) Dispatch(c actor.Context, msg camel.Message) bool {
	// no children
	if v.pids.Empty() {
		return true
	}

	sender := c.Sender()

	if !v.pids.Contains(sender) {
		// this is not a message coming from a children, send to the first one
		c.Send(v.pids.Get(0), msg)

		return false
	}

	for i := 0; i < v.pids.Len(); i++ {
		pid := v.pids.Get(i)

		if pid.Equal(sender) && i != v.pids.Len()-1 {
			// send the message to the next one
			c.Send(v.pids.Get(i+1), msg)
			return false
		}
	}

	return true
}
