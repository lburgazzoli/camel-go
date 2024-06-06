package choice

import (
	"context"
	"testing"
	"time"

	"github.com/lburgazzoli/camel-go/pkg/core/processors/choice/otherwise"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/choice/when"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/language"
	"github.com/lburgazzoli/camel-go/pkg/core/language/jq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
)

func TestChoiceWhen(t *testing.T) {
	g := support.With(t)
	c := camel.ExtractContext(g.Ctx())

	choice := New()
	choice.When = []*when.When{
		when.New(
			when.WithExpression(language.Language{
				Jq: &jq.Jq{Definition: jq.Definition{Expression: `.foo == "bar"`}},
			}),
			when.WithProcessor(func(ctx context.Context, m camel.Message) error {
				m.SetContent("branch: bar")
				return nil
			}),
		),
		when.New(
			when.WithExpression(language.Language{
				Jq: &jq.Jq{Definition: jq.Definition{Expression: `.foo == "baz"`}},
			}),
			when.WithProcessor(func(ctx context.Context, m camel.Message) error {
				m.SetContent("branch: baz")
				return nil
			}),
		),
	}

	chv, err := choice.Reify(g.Ctx())

	require.NoError(t, err)
	assert.NotNil(t, chv)

	h := support.NewRootVerticle(chv)

	chp, err := c.Spawn(h)
	require.NoError(t, err)
	require.NotNil(t, chp)

	{
		msg := c.NewMessage()
		msg.SetContent(`{ "foo": "bar" }`)
		require.NoError(t, err)

		err = c.SendTo(chp, msg)
		require.NoError(t, err)

		res, err := h.Get(1 * time.Minute)
		require.NoError(t, err)
		require.Equal(t, "branch: bar", res.Content())
	}

	{
		msg := c.NewMessage()
		msg.SetContent(`{ "foo": "baz" }`)
		require.NoError(t, err)

		err = c.SendTo(chp, msg)
		require.NoError(t, err)

		res, err := h.Get(1 * time.Minute)
		require.NoError(t, err)
		require.Equal(t, "branch: baz", res.Content())
	}
}

func TestChoiceOtherWise(t *testing.T) {
	g := support.With(t)
	c := camel.ExtractContext(g.Ctx())

	choice := New()
	choice.When = []*when.When{
		when.New(
			when.WithExpression(language.Language{
				Jq: &jq.Jq{Definition: jq.Definition{Expression: `.foo == "bar"`}},
			}),
			when.WithProcessor(func(ctx context.Context, m camel.Message) error {
				m.SetContent("branch: bar")
				return nil
			}),
		),
		when.New(
			when.WithExpression(language.Language{
				Jq: &jq.Jq{Definition: jq.Definition{Expression: `.foo == "baz"`}},
			}),
			when.WithProcessor(func(ctx context.Context, m camel.Message) error {
				m.SetContent("branch: baz")
				return nil
			}),
		),
	}
	choice.Otherwise = otherwise.New(
		otherwise.WithProcessor(func(ctx context.Context, m camel.Message) error {
			m.SetContent("branch: otherwise")
			return nil
		}),
	)

	chv, err := choice.Reify(g.Ctx())

	require.NoError(t, err)
	assert.NotNil(t, chv)

	h := support.NewRootVerticle(chv)

	chp, err := c.Spawn(h)
	require.NoError(t, err)
	require.NotNil(t, chp)

	{
		msg := c.NewMessage()
		msg.SetContent(`{ "foo": "bar" }`)
		require.NoError(t, err)

		err = c.SendTo(chp, msg)
		require.NoError(t, err)

		res, err := h.Get(1 * time.Minute)
		require.NoError(t, err)
		require.Equal(t, "branch: bar", res.Content())
	}

	{
		msg := c.NewMessage()
		msg.SetContent(`{ "foo": "baz" }`)
		require.NoError(t, err)

		err = c.SendTo(chp, msg)
		require.NoError(t, err)

		res, err := h.Get(1 * time.Minute)
		require.NoError(t, err)
		require.Equal(t, "branch: baz", res.Content())
	}

	{
		msg := c.NewMessage()
		msg.SetContent(`{ "bar": "baz" }`)
		require.NoError(t, err)

		err = c.SendTo(chp, msg)
		require.NoError(t, err)

		res, err := h.Get(1 * time.Minute)
		require.NoError(t, err)
		require.Equal(t, "branch: otherwise", res.Content())
	}
}
