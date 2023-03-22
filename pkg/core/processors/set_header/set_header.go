// //go:build steps_process || steps_all

package process

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

const TAG = "setHeader"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &Process{
			DefaultVerticle: processors.NewDefaultVerticle(),
		}
	}
}

type Process struct {
	processors.DefaultVerticle `yaml:",inline"`

	Name     string `yaml:"name"`
	Language `yaml:",inline"`
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

func (p *Process) Reify(ctx context.Context) (camel.Verticle, error) {
	camelContext := camel.GetContext(ctx)

	if p.Name == "" {
		return nil, camelerrors.MissingParameterf("name", "failure processing %s", TAG)
	}
	if p.Constant == nil {
		return nil, camelerrors.MissingParameterf("constant", "failure processing %s", TAG)
	}
	if p.Constant.Value == "" {
		return nil, camelerrors.MissingParameterf("constant.value", "failure processing %s", TAG)
	}

	p.SetContext(camelContext)

	return p, nil
}

func (p *Process) Receive(ac actor.Context) {
	msg, ok := ac.Message().(camel.Message)
	if ok {
		_ = msg.SetExtension(p.Name, p.Constant.Value)

		p.Dispatch(ac, msg)
	}
}
