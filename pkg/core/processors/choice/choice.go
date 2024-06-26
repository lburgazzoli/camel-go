package choice

import (
	"context"

	"github.com/lburgazzoli/camel-go/pkg/core/processors/choice/otherwise"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/choice/when"

	"github.com/lburgazzoli/camel-go/pkg/util/uuid"

	"github.com/lburgazzoli/camel-go/pkg/core/verticles"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/pkg/errors"
)

const TAG = "choice"

func init() {
	processors.Types[TAG] = func() interface{} {
		return New()
	}
}

func New() *Choice {
	return &Choice{
		DefaultVerticle: processors.NewDefaultVerticle(),
	}
}

type Choice struct {
	processors.DefaultVerticle `yaml:",inline"`

	When      []*when.When         `yaml:"when,omitempty"`
	Otherwise *otherwise.Otherwise `yaml:"otherwise,omitempty"`
}

func (c *Choice) Reify(ctx context.Context) (camel.Verticle, error) {
	camelContext := camel.ExtractContext(ctx)

	c.SetContext(camelContext)

	return c, nil
}

func (c *Choice) Receive(ac actor.Context) {
	ctx := verticles.NewContext(c.Context(), ac)

	switch msg := ac.Message().(type) {
	case *actor.Started:
		c.onStarted(ctx, ac, msg)
	case camel.Message:
		c.onMessage(ctx, ac, msg)
	case processors.StepsDone:
		c.onDone(ctx, ac, msg.M)
	}
}

func (c *Choice) onStarted(ctx context.Context, ac actor.Context, _ *actor.Started) {
	for w := range c.When {
		if c.When[w].Identity == "" {
			c.When[w].Identity = uuid.New()
		}

		v, err := c.When[w].Reify(ctx)
		if err != nil {
			panic(errors.Wrapf(err, "unable to reify verticle with id <%s>", c.When[w].ID()))
		}

		p, err := verticles.Spawn(ac, v)
		if err != nil {
			panic(errors.Wrapf(err, "unable to spawn verticle with id <%s>", c.When[w].ID()))
		}

		// Horrible
		c.When[w].PID = p
	}

	if c.Otherwise != nil {
		if c.Otherwise.Identity == "" {
			c.Otherwise.Identity = uuid.New()
		}

		v, err := c.Otherwise.Reify(ctx)
		if err != nil {
			panic(errors.Wrapf(err, "unable to reify verticle with id <%s>", c.Otherwise.ID()))
		}

		p, err := verticles.Spawn(ac, v)
		if err != nil {
			panic(errors.Wrapf(err, "unable to spawn verticle with id <%s>", c.Otherwise.ID()))
		}

		// Horrible
		c.Otherwise.PID = p
	}
}

func (c *Choice) onMessage(ctx context.Context, ac actor.Context, msg camel.Message) {
	var matches bool

	for i := range c.When {
		m, err := c.When[i].Matches(ctx, msg)
		if err != nil {
			panic(err)
		}

		matches = m

		if matches {
			ac.Request(c.When[i].PID, msg)
			break
		}
	}

	if !matches && c.Otherwise != nil {
		ac.Request(c.Otherwise.PID, msg)
	}
}

func (c *Choice) onDone(_ context.Context, ac actor.Context, msg camel.Message) {
	// all done, unwrap and send to parent
	if ac.Parent() != nil {
		ac.Request(ac.Parent(), msg)
	}
}
