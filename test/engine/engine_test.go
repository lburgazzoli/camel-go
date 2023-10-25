// //go:build components_all || components_timer || steps_all || steps_process

package engine

import (
	"context"
	"strings"

	"github.com/stretchr/testify/require"

	"testing"
	"time"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
	// helper to include everything.
	_ "github.com/lburgazzoli/camel-go/pkg/components/dapr/pubsub"
	_ "github.com/lburgazzoli/camel-go/pkg/components/http"
	_ "github.com/lburgazzoli/camel-go/pkg/components/kafka"
	_ "github.com/lburgazzoli/camel-go/pkg/components/log"
	_ "github.com/lburgazzoli/camel-go/pkg/components/mqtt/v3"
	_ "github.com/lburgazzoli/camel-go/pkg/components/wasm"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/choice"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/process"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/to"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/transform"

	"github.com/lburgazzoli/camel-go/pkg/components/timer"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/stretchr/testify/assert"
)

const simpleYAML = `
- route:
    from:
      uri: "timer:foo"
      steps:
        - process:
            ref: "consumer-1"
        - process:
            ref: "consumer-2"
`

func TestSimpleYAML(t *testing.T) {
	g := support.With(t)
	c := camel.ExtractContext(g.Ctx())

	content := uuid.New()
	wg := make(chan camel.Message)

	c.Registry().Set("consumer-1", func(_ context.Context, message camel.Message) error {
		message.SetContent(content)
		return nil
	})
	c.Registry().Set("consumer-2", func(_ context.Context, message camel.Message) error {
		wg <- message
		return nil
	})

	err := c.LoadRoutes(g.Ctx(), strings.NewReader(simpleYAML))
	require.NoError(t, err)

	select {
	case msg := <-wg:
		a, ok := msg.Attribute(timer.AttributeTimerFiredCount)
		assert.True(t, ok)
		assert.Equal(t, "1", a)
		assert.Equal(t, content, msg.Content())
	case <-time.After(60 * time.Second):
		assert.Fail(t, "timeout")
	}
}

const simpleWASM = `
- route:
    from:
      uri: "timer:foo"
      steps:
        - process:
            ref: "consumer-1"
        - transform:
            wasm: 
              path: "../../etc/wasm/fn/simple_process.wasm"
        - process:
            ref: "consumer-2"
`

func TestSimpleWASM(t *testing.T) {
	g := support.With(t)
	c := camel.ExtractContext(g.Ctx())

	wg := make(chan camel.Message)

	c.Registry().Set("consumer-1", func(_ context.Context, message camel.Message) error {
		message.SetSubject("consumer-1")
		return nil
	})
	c.Registry().Set("consumer-2", func(_ context.Context, message camel.Message) error {
		wg <- message
		return nil
	})

	err := c.LoadRoutes(g.Ctx(), strings.NewReader(simpleWASM))
	require.NoError(t, err)

	select {
	case msg := <-wg:
		assert.Equal(t, "consumer-1", msg.Subject())

		c, ok := msg.Content().([]byte)
		assert.True(t, ok)
		assert.Equal(t, "hello from wasm", string(c))

	case <-time.After(5 * time.Second):
		assert.Fail(t, "timeout")
	}
}

const simpleInlineWASM = `
- route:
    from:
      uri: "timer:foo"
      steps:
        - process:
            ref: "consumer-1"
        - transform:
            wasm: "../../etc/wasm/fn/simple_process.wasm"
        - process:
            ref: "consumer-2"
`

func TestSimpleInlineWASM(t *testing.T) {
	g := support.With(t)
	c := camel.ExtractContext(g.Ctx())
	wg := make(chan camel.Message)

	c.Registry().Set("consumer-1", func(_ context.Context, message camel.Message) error {
		message.SetSubject("consumer-1")
		return nil
	})
	c.Registry().Set("consumer-2", func(_ context.Context, message camel.Message) error {
		wg <- message
		return nil
	})

	err := c.LoadRoutes(g.Ctx(), strings.NewReader(simpleInlineWASM))
	require.NoError(t, err)

	select {
	case msg := <-wg:
		assert.Equal(t, "consumer-1", msg.Subject())

		c, ok := msg.Content().([]byte)
		assert.True(t, ok)
		assert.Equal(t, "hello from wasm", string(c))

	case <-time.After(5 * time.Second):
		assert.Fail(t, "timeout")
	}
}

const simpleInlineImageWASM = `
- route:
    from:
      uri: "timer:foo"
      steps:
        - process:
            ref: "consumer-1"
        - transform:
            wasm: "quay.io/lburgazzoli/camel-go-wasm?etc/wasm/fn/simple_process.wasm"
        - process:
            ref: "consumer-2"
`

func TestSimpleInlineImageWASM(t *testing.T) {
	g := support.With(t)
	c := camel.ExtractContext(g.Ctx())

	wg := make(chan camel.Message)

	c.Registry().Set("consumer-1", func(_ context.Context, message camel.Message) error {
		message.SetSubject("consumer-1")
		return nil
	})
	c.Registry().Set("consumer-2", func(_ context.Context, message camel.Message) error {
		wg <- message
		return nil
	})

	err := c.LoadRoutes(g.Ctx(), strings.NewReader(simpleInlineImageWASM))
	require.NoError(t, err)

	select {
	case msg := <-wg:
		assert.Equal(t, "consumer-1", msg.Subject())

		c, ok := msg.Content().([]byte)
		assert.True(t, ok)
		assert.Equal(t, "hello from wasm", string(c))

	case <-time.After(5 * time.Second):
		assert.Fail(t, "timeout")
	}
}
