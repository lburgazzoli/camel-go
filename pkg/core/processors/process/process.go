// //go:build steps_process || steps_all

package process

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/core/registry"
)

const TAG = "process"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &Process{
			DefaultVerticle: processors.NewDefaultVerticle(),
		}
	}
}

type Process struct {
	processors.DefaultVerticle `yaml:",inline"`

	Ref       string `yaml:"ref"`
	processor camel.Processor
}

func (p *Process) Reify(_ context.Context, camelContext camel.Context) (string, error) {

	if p.Ref == "" {
		return "", camelerrors.MissingParameterf("ref", "failure processing %s", TAG)
	}

	proc, ok := registry.GetAs[camel.Processor](camelContext.Registry(), p.Ref)
	if !ok {
		return "", camelerrors.MissingParameterf("ref", "failure processing %s", TAG)
	}

	p.SetContext(camelContext)
	p.processor = proc

	return p.Identity, camelContext.Spawn(p)
}

func (p *Process) Receive(c actor.Context) {
	msg, ok := c.Message().(camel.Message)
	if ok {
		if err := p.processor(msg); err != nil {
			panic(err)
		}

		p.Dispatch(msg)
	}
}
