//go:build steps_process || steps_all

package process

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

const TAG = "process"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &Process{}
	}
}

type Process struct {
	outputs *actor.PIDSet

	Ref string `yaml:"ref"`
}

func (p *Process) Next(pid *actor.PID) {
	if p.outputs == nil {
		p.outputs = actor.NewPIDSet()
	}

	p.outputs.Add(pid)
}

func (p *Process) Reify(ctx api.Context) (*actor.PID, error) {
	if p.Ref == "" {
		return nil, camelerrors.MissingParameterf("ref", "failure processing %s", TAG)
	}

	val, ok := ctx.Registry().Get(p.Ref)
	if !ok {
		return nil, camelerrors.MissingParameterf("ref", "failure processing %s", TAG)
	}

	proc, ok := val.(func(api.Message))
	if !ok {
		return nil, camelerrors.InvalidParameterf("ref", "failure retrieving value from registry")
	}

	pid := ctx.Spawn(&processActor{
		camelContext:   ctx,
		camelProcessor: proc,
		outputs:        p.outputs,
	})

	return pid, nil
}

type processActor struct {
	camelContext   api.Context
	camelProcessor api.Processor
	outputs        *actor.PIDSet
}

func (p *processActor) Receive(c actor.Context) {
	switch msg := c.Message().(type) {
	case api.Message:
		p.camelProcessor(msg)

		if p.outputs != nil {
			p.outputs.ForEach(func(_ int, pid *actor.PID) {
				c.Send(pid, msg)
			})
		}
	}
}
