package choice

import (
	"context"
	"testing"
	"time"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/language"
	"github.com/lburgazzoli/camel-go/pkg/core/language/jq"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
)

func TestChoice(t *testing.T) {
	support.Run(t, "when", func(t *testing.T, ctx context.Context) {
		t.Helper()

		c := camel.ExtractContext(ctx)

		choice := New()
		choice.When = []*When{
			NewWhen(
				language.Language{
					Jq: &jq.Jq{Definition: jq.Definition{Expression: `.foo == "bar"`}},
				},
				processors.NewStep(support.NewProcessorsVerticle(func(ctx context.Context, m camel.Message) error {
					m.SetContent("branch: bar")
					return nil
				})),
			),
			NewWhen(
				language.Language{
					Jq: &jq.Jq{Definition: jq.Definition{Expression: `.foo == "baz"`}},
				},
				processors.NewStep(support.NewProcessorsVerticle(func(ctx context.Context, m camel.Message) error {
					m.SetContent("branch: baz")
					return nil
				})),
			),
		}

		chv, err := choice.Reify(ctx)

		assert.Nil(t, err)
		assert.NotNil(t, chv)

		h := support.NewRootVerticle(chv)

		chp, err := c.Spawn(h)
		require.Nil(t, err)
		require.NotNil(t, chp)

		{
			msg, err := message.New()
			msg.SetContent(`{ "foo": "bar" }`)
			assert.Nil(t, err)

			err = c.SendTo(chp, msg)
			require.Nil(t, err)

			res, err := h.Get(1 * time.Minute)
			require.Nil(t, err)
			require.Equal(t, "branch: bar", res.Content())
		}

		{
			msg, err := message.New()
			msg.SetContent(`{ "foo": "baz" }`)
			assert.Nil(t, err)

			err = c.SendTo(chp, msg)
			require.Nil(t, err)

			res, err := h.Get(1 * time.Minute)
			require.Nil(t, err)
			require.Equal(t, "branch: baz", res.Content())
		}
	})

	support.Run(t, "otherwise", func(t *testing.T, ctx context.Context) {
		t.Helper()

		c := camel.ExtractContext(ctx)

		choice := New()
		choice.When = []*When{
			NewWhen(
				language.Language{
					Jq: &jq.Jq{Definition: jq.Definition{Expression: `.foo == "bar"`}},
				},
				processors.NewStep(support.NewProcessorsVerticle(func(ctx context.Context, m camel.Message) error {
					m.SetContent("branch: bar")
					return nil
				})),
			),
			NewWhen(
				language.Language{
					Jq: &jq.Jq{Definition: jq.Definition{Expression: `.foo == "baz"`}},
				},
				processors.NewStep(support.NewProcessorsVerticle(func(ctx context.Context, m camel.Message) error {
					m.SetContent("branch: baz")
					return nil
				})),
			),
		}
		choice.Otherwise = NewOtherwise(
			processors.NewStep(support.NewProcessorsVerticle(func(ctx context.Context, m camel.Message) error {
				m.SetContent("branch: otherwise")
				return nil
			})),
		)

		chv, err := choice.Reify(ctx)

		assert.Nil(t, err)
		assert.NotNil(t, chv)

		h := support.NewRootVerticle(chv)

		chp, err := c.Spawn(h)
		require.Nil(t, err)
		require.NotNil(t, chp)

		{
			msg, err := message.New()
			msg.SetContent(`{ "foo": "bar" }`)
			assert.Nil(t, err)

			err = c.SendTo(chp, msg)
			require.Nil(t, err)

			res, err := h.Get(1 * time.Minute)
			require.Nil(t, err)
			require.Equal(t, "branch: bar", res.Content())
		}

		{
			msg, err := message.New()
			msg.SetContent(`{ "foo": "baz" }`)
			assert.Nil(t, err)

			err = c.SendTo(chp, msg)
			require.Nil(t, err)

			res, err := h.Get(1 * time.Minute)
			require.Nil(t, err)
			require.Equal(t, "branch: baz", res.Content())
		}

		{
			msg, err := message.New()
			msg.SetContent(`{ "bar": "baz" }`)
			assert.Nil(t, err)

			err = c.SendTo(chp, msg)
			require.Nil(t, err)

			res, err := h.Get(1 * time.Minute)
			require.Nil(t, err)
			require.Equal(t, "branch: otherwise", res.Content())
		}
	})
}
