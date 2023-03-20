package choice

import (
	"context"
	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/pkg/errors"
)

const TAG = "choice"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &Choice{
			DefaultVerticle: processors.NewDefaultVerticle(),
		}
	}
}

type Choice struct {
	processors.DefaultVerticle `yaml:",inline"`

	When []When `yaml:"when,omitempty"`
}

func (c *Choice) Reify(_ context.Context, camelContext camel.Context) (string, error) {
	c.SetContext(camelContext)

	var last string

	for w := range c.When {

		for i := len(c.When[w].Steps) - 1; i >= 0; i-- {
			if last != "" {
				f.Steps[i].Next(last)
			}

			pid, err := f.Steps[i].Reify(ctx, camelContext)
			if err != nil {
				return "", errors.Wrapf(err, "error creating step")
			}

			last = pid
		}

		if last != "" {
			f.Endpoint.Next(last)
		}
	}

	return "", camelerrors.NotImplemented("TODO")
}

func (c *Choice) Receive(ac actor.Context) {
	msg, ok := ac.Message().(camel.Message)
	if ok {

		ctx := context.Background()

		for i := range c.When {
			matches, err := c.When[i].Matches(ctx, msg)
			if err != nil {
				panic(err)
			}

			if matches {
				c.When[i].Dispatch(msg)
				break
			}
		}
	}
}
