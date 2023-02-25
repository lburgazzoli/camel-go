////go:build steps_process || steps_all

package process

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/core/registry"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

const TAG = "process"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &Process{}
	}
}

type Process struct {
	api.Identifiable
	api.WithOutputs

	Identity string `yaml:"id"`
	Ref      string `yaml:"ref"`
}

func (p *Process) ID() string {
	return p.Identity
}

func (p *Process) Reify(ctx api.Context) (*actor.PID, error) {

	if p.Ref == "" {
		return nil, camelerrors.MissingParameterf("ref", "failure processing %s", TAG)
	}

	proc, ok := registry.GetAs[api.Processor](ctx.Registry(), p.Ref)
	if !ok {
		return nil, camelerrors.MissingParameterf("ref", "failure processing %s", TAG)
	}

	id := p.Identity
	if id == "" {
		id = uuid.New()
	}

	return ctx.Spawn(id, &processActor{
		context:   ctx,
		processor: proc,
		outputs:   p.Outputs(),
	})
}

type processActor struct {
	context   api.Context
	processor api.Processor
	outputs   []*actor.PID
}

func (p *processActor) Receive(c actor.Context) {
	msg, ok := c.Message().(api.Message)
	if ok {
		p.processor(msg)

		for i := range p.outputs {
			c.Send(p.outputs[i], msg)
		}
	}
}
