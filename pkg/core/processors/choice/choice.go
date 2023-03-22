package choice

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
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

	When      []*When    `yaml:"when,omitempty"`
	Otherwise *Otherwise `yaml:"otherwise,omitempty"`
}

func (c *Choice) Reify(ctx context.Context) (string, error) {
	camelContext := camel.GetContext(ctx)

	c.SetContext(camelContext)

	for w := range c.When {
		when := c.When[w]

		if err := when.Configure(ctx, camelContext); err != nil {
			return "", errors.Wrapf(err, "error configuring when %s", when.ID())
		}

		if err := processors.ReifySteps(ctx, when, when.Steps); err != nil {
			return "", errors.Wrapf(err, "error creating when steps")
		}
	}

	if c.Otherwise != nil {
		if err := c.Otherwise.Configure(ctx, camelContext); err != nil {
			return "", errors.Wrapf(err, "error configuring otherwhise %s", c.Otherwise.ID())
		}

		if err := processors.ReifySteps(ctx, c.Otherwise, c.Otherwise.Steps); err != nil {
			return "", errors.Wrapf(err, "error creating otherwhise steps")
		}
	}

	return c.Identity, camelContext.Spawn(c)
}

func (c *Choice) Receive(ac actor.Context) {
	msg, ok := ac.Message().(camel.Message)
	if ok {

		ctx := context.Background()

		var matches bool
		var err error

		for i := range c.When {
			matches, err = c.When[i].Matches(ctx, msg)
			if err != nil {
				panic(err)
			}

			if matches {
				c.When[i].Dispatch(msg)
				break
			}
		}

		if !matches && c.Otherwise != nil {
			c.Otherwise.Dispatch(msg)
		}
	}
}
