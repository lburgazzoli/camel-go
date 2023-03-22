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
		return NewProcess()
	}
}

func NewProcess() *Process {
	return &Process{
		DefaultVerticle: processors.NewDefaultVerticle(),
	}
}

func NewProcessWithRef(ref string) *Process {
	return &Process{
		DefaultVerticle: processors.NewDefaultVerticle(),
		Ref:             ref,
	}
}

type Process struct {
	processors.DefaultVerticle `yaml:",inline"`

	Ref       string `yaml:"ref"`
	processor camel.Processor
}

func (p *Process) Reify(ctx context.Context) (camel.Verticle, error) {
	camelContext := camel.GetContext(ctx)

	if p.Ref == "" {
		return nil, camelerrors.MissingParameterf("ref", "failure processing %s", TAG)
	}

	proc, ok := registry.GetAs[camel.Processor](camelContext.Registry(), p.Ref)
	if !ok {
		return nil, camelerrors.MissingParameterf("ref", "unable to lookup processor %s from registry", p.Ref)
	}

	p.SetContext(camelContext)
	p.processor = proc

	return p, nil
}

func (p *Process) Receive(c actor.Context) {
	msg, ok := c.Message().(camel.Message)
	if ok {
		if err := p.processor(context.Background(), msg); err != nil {
			panic(err)
		}

		p.Dispatch(c, msg)
	}
}
