////go:build steps_process || steps_all

package process

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/core/registry"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

const TAG = "process"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &Process{
			Identity: uuid.New(),
		}
	}
}

type Process struct {
	camel.Identifiable
	camel.WithOutputs

	Identity string `yaml:"id"`
	Ref      string `yaml:"ref"`

	context   camel.Context
	processor camel.Processor
}

func (p *Process) ID() string {
	return p.Identity
}

func (p *Process) Reify(_ context.Context, camelContext camel.Context) (string, error) {

	if p.Ref == "" {
		return "", camelerrors.MissingParameterf("ref", "failure processing %s", TAG)
	}

	proc, ok := registry.GetAs[camel.Processor](camelContext.Registry(), p.Ref)
	if !ok {
		return "", camelerrors.MissingParameterf("ref", "failure processing %s", TAG)
	}

	p.context = camelContext
	p.processor = proc

	return p.Identity, camelContext.Spawn(p)
}

func (p *Process) Receive(c actor.Context) {
	msg, ok := c.Message().(camel.Message)
	if ok {
		p.processor(msg)

		for _, pid := range p.Outputs() {
			if err := p.context.Send(pid, msg); err != nil {
				panic(err)
			}
		}
	}
}
