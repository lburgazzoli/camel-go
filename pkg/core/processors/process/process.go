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
		return New()
	}
}

func New(opts ...OptionFn) *Process {
	answer := &Process{
		DefaultVerticle: processors.NewDefaultVerticle(),
	}

	for _, o := range opts {
		o(answer)
	}

	return answer
}

type Process struct {
	processors.DefaultVerticle `yaml:",inline"`

	Ref       string `yaml:"ref"`
	processor camel.Processor
}

func (p *Process) Reify(ctx context.Context) (camel.Verticle, error) {
	camelContext := camel.ExtractContext(ctx)

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

func (p *Process) Receive(ac actor.Context) {
	msg, ok := ac.Message().(camel.Message)
	if ok {
		ctx := camel.Wrap(context.Background(), p.Context())

		if err := p.processor(ctx, msg); err != nil {
			panic(err)
		}

		ac.Request(ac.Sender(), msg)
	}
}
