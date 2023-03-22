package wasm

import (
	"context"
	"strings"

	"testing"
	"time"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/stretchr/testify/assert"

	// test support.
	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"

	// enable components.
	_ "github.com/lburgazzoli/camel-go/pkg/components/log"
	_ "github.com/lburgazzoli/camel-go/pkg/components/timer"

	// enable processors.
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/process"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/to"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/transform"
)

const simpleComponentWASM = `
- route:
    from:
      uri: "timer:foo"
      steps:
        - process:
            ref: "consumer-1"
        - to:
            uri: "wasm:../../../etc/wasm/fn/simple_logger.wasm"
        - process:
            ref: "consumer-2"
`

func TestSimpleComponentWASM(t *testing.T) {
	support.Run(t, "run", func(t *testing.T, ctx context.Context) {
		t.Helper()

		wg := make(chan camel.Message)

		c := camel.GetContext(ctx)

		c.Registry().Set("consumer-1", func(_ context.Context, message camel.Message) error {
			message.SetContent("consumer-1")
			return nil
		})
		c.Registry().Set("consumer-2", func(_ context.Context, message camel.Message) error {
			wg <- message
			return nil
		})

		err := c.LoadRoutes(ctx, strings.NewReader(simpleComponentWASM))
		assert.Nil(t, err)

		select {
		case msg := <-wg:
			c, ok := msg.Content().(string)
			assert.True(t, ok)
			assert.Equal(t, "consumer-1", c)

		case <-time.After(5 * time.Second):
			assert.Fail(t, "timeout")
		}
	})
}

const simpleComponentImageWASM = `
- route:
    from:
      uri: "timer:foo"
      steps:
        - process:
            ref: "consumer-1"
        - to:
            uri: "wasm:etc/wasm/fn/simple_logger.wasm?image=docker.io/lburgazzoli/camel-go:latest"
        - process:
            ref: "consumer-2"
`

func TestSimpleComponentImageWASM(t *testing.T) {
	support.Run(t, "run", func(t *testing.T, ctx context.Context) {
		t.Helper()

		wg := make(chan camel.Message)

		c := camel.GetContext(ctx)

		c.Registry().Set("consumer-1", func(_ context.Context, message camel.Message) error {
			message.SetContent("consumer-1")
			return nil
		})
		c.Registry().Set("consumer-2", func(_ context.Context, message camel.Message) error {
			wg <- message
			return nil
		})

		err := c.LoadRoutes(ctx, strings.NewReader(simpleComponentImageWASM))
		assert.Nil(t, err)

		select {
		case msg := <-wg:
			c, ok := msg.Content().(string)
			assert.True(t, ok)
			assert.Equal(t, "consumer-1", c)

		case <-time.After(5 * time.Second):
			assert.Fail(t, "timeout")
		}
	})
}
