////go:build steps_process || steps_all

package process

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

const TAG = "setHeader"

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
	Name     string `yaml:"name"`
	Language `yaml:",inline"`

	context   camel.Context
	processor camel.Processor
}

type Language struct {
	Constant *LanguageConstant `yaml:"constant,omitempty"`
}

type LanguageConstant struct {
	Value string `yaml:"value"`
}

func (p *Process) ID() string {
	return p.Identity
}

func (p *Process) Reify(_ context.Context, camelContext camel.Context) (string, error) {
	if p.Name == "" {
		return "", camelerrors.MissingParameterf("name", "failure processing %s", TAG)
	}
	if p.Constant == nil {
		return "", camelerrors.MissingParameterf("constant", "failure processing %s", TAG)
	}
	if p.Constant.Value == "" {
		return "", camelerrors.MissingParameterf("constant.value", "failure processing %s", TAG)
	}

	p.context = camelContext

	return p.Identity, camelContext.Spawn(p)
}

func (p *Process) Receive(c actor.Context) {
	msg, ok := c.Message().(camel.Message)
	if ok {
		msg.SetExtension(p.Name, p.Constant.Value)

		for _, pid := range p.Outputs() {
			if err := p.context.Send(pid, msg); err != nil {
				panic(err)
			}
		}
	}
}
