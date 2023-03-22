package choice

import (
	"context"

	"github.com/lburgazzoli/camel-go/pkg/core/verticles"

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

func (c *Choice) Reify(ctx context.Context) (camel.Verticle, error) {
	camelContext := camel.GetContext(ctx)

	c.SetContext(camelContext)

	return c, nil
}

func (c *Choice) Receive(ac actor.Context) {
	ctx := context.Background()

	switch msg := ac.Message().(type) {
	case *actor.Started:
		ctx := verticles.NewContext(c.Context(), ac)

		for w := range c.When {
			v, err := c.When[w].Reify(ctx)
			if err != nil {
				panic(errors.Wrapf(err, "unable to reify verticle with id %s", c.When[w].ID()))
			}

			pid, err := verticles.Spawn(ac, v)
			if err != nil {
				panic(errors.Wrapf(err, "unable to spawn verticle with id %s", c.When[w].ID()))
			}

			// Ugly, very ugly
			c.When[w].pid = pid
		}

		if c.Otherwise != nil {
			v, err := c.Otherwise.Reify(ctx)
			if err != nil {
				panic(errors.Wrapf(err, "unable to reify verticle with id %s", c.Otherwise.ID()))
			}

			pid, err := verticles.Spawn(ac, v)
			if err != nil {
				panic(errors.Wrapf(err, "unable to spawn verticle with id %s", c.Otherwise.ID()))
			}

			// Ugly, very ugly
			c.Otherwise.pid = pid
		}
	case camel.Message:
		var matches bool
		var err error

		for i := range c.When {
			matches, err = c.When[i].Matches(ctx, msg)
			if err != nil {
				panic(err)
			}

			if matches {
				ac.Send(c.When[i].pid, msg)
				break
			}
		}

		if !matches && c.Otherwise != nil {
			ac.Send(c.Otherwise.pid, msg)
		}
	}
}
