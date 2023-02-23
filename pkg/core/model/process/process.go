package process

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/core/model"
)

const TAG = "process"

func init() {
	model.Types[TAG] = func() interface{} {
		return &Process{}
	}
}

type Process struct {
	Ref    string       `yaml:"ref"`
	SendTo []*actor.PID `yaml:"-"`
}

func (e *Process) Reify(ctx api.Context) (actor.Actor, error) {
	if e.Ref == "" {
		return nil, errors.MissingParameterf("ref", "failure processing %s", TAG)
	}

	val, ok := ctx.Registry().Get(e.Ref)
	if !ok {
		return nil, errors.MissingParameterf("ref", "failure processing %s", TAG)
	}

	p, ok := val.(func(api.Message))
	if !ok {
		return nil, errors.InvalidParameterf("ref", "failure retrieving value from registry")
	}

	a := processActor{
		camelContext:   ctx,
		camelProcessor: p,
		targets:        e.SendTo,
	}

	return &a, nil
}

type processActor struct {
	camelContext   api.Context
	camelProcessor api.Processor
	targets        []*actor.PID
}

func (p *processActor) Receive(c actor.Context) {
	switch msg := c.Message().(type) {
	case *actor.Started:
		break
	case *actor.Stopping:
		break
	case *actor.Stopped:
		break
	case *actor.Restarting:
		break
	case api.Message:
		p.camelProcessor(msg)
		
		for i := range p.targets {
			c.Send(p.targets[i], msg)
		}
	}
}
