// //go:build steps_process || steps_all

package setbody

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

const TAG = "setBody"

func init() {
	processors.Types[TAG] = func() interface{} {
		return New()
	}
}

func New() *SetBody {
	return &SetBody{
		DefaultVerticle: processors.NewDefaultVerticle(),
	}
}

type SetBody struct {
	processors.DefaultVerticle `yaml:",inline"`
	Language                   `yaml:",inline"`
}

type Language struct {
	Constant *LanguageConstant `yaml:"constant,omitempty"`
}

type LanguageConstant struct {
	Value string `yaml:"value"`
}

func (p *SetBody) Reify(ctx context.Context) (camel.Verticle, error) {
	camelContext := camel.ExtractContext(ctx)

	if p.Constant == nil {
		return nil, camelerrors.MissingParameterf("constant", "failure processing %s", TAG)
	}
	if p.Constant.Value == "" {
		return nil, camelerrors.MissingParameterf("constant.value", "failure processing %s", TAG)
	}

	p.SetContext(camelContext)

	return p, nil
}

func (p *SetBody) Receive(ac actor.Context) {
	msg, ok := ac.Message().(camel.Message)
	if ok {
		msg.SetContent(p.Constant.Value)

		ac.Request(ac.Sender(), msg)
	}
}
