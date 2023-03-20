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

func (c *Choice) Reify(ctx context.Context, camelContext camel.Context) (string, error) {
	c.SetContext(camelContext)

	for w := range c.When {
		var last string

		when := c.When[w]

		if err := when.Configure(ctx, camelContext); err != nil {
			return "", errors.Wrapf(err, "error configuring when %s", when.ID())
		}

		for s := len(when.Steps) - 1; s >= 0; s-- {
			step := when.Steps[s]

			if last != "" {
				step.Next(last)
			}

			id, err := step.Reify(ctx, camelContext)
			if err != nil {
				return "", errors.Wrapf(err, "error creating when step")
			}

			last = id
		}

		if last != "" {
			when.Next(last)
		}
	}

	if c.Otherwise != nil {
		if err := c.Otherwise.Configure(ctx, camelContext); err != nil {
			return "", errors.Wrapf(err, "error configuring otherwhise %s", c.Otherwise.ID())
		}

		var last string

		for s := len(c.Otherwise.Steps) - 1; s >= 0; s-- {
			step := c.Otherwise.Steps[s]

			if last != "" {
				step.Next(last)
			}

			id, err := step.Reify(ctx, camelContext)
			if err != nil {
				return "", errors.Wrapf(err, "error creating otherwhise step")
			}

			last = id
		}

		if last != "" {
			c.Otherwise.Next(last)
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
