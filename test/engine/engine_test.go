// //go:build components_all || components_timer || steps_all || steps_process

package engine

import (
	"context"
	"strings"

	"github.com/lburgazzoli/camel-go/pkg/core"
	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"testing"
	"time"

	_ "github.com/lburgazzoli/camel-go/pkg/components/kafka"
	_ "github.com/lburgazzoli/camel-go/pkg/components/log"
	_ "github.com/lburgazzoli/camel-go/pkg/components/mqtt"
	"github.com/lburgazzoli/camel-go/pkg/components/timer"
	_ "github.com/lburgazzoli/camel-go/pkg/components/wasm"

	"github.com/lburgazzoli/camel-go/pkg/util/uuid"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelc "github.com/lburgazzoli/camel-go/pkg/core/context"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/process"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/to"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/transform"
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
	support.Run(t, "run", func(t *testing.T, ctx context.Context) {
		t.Helper()

		content := uuid.New()
		wg := make(chan camel.Message)

		c := camel.GetContext(ctx)

		c.Registry().Set("consumer-1", func(_ context.Context, message camel.Message) error {
			message.SetContent(content)
			return nil
		})
		c.Registry().Set("consumer-2", func(_ context.Context, message camel.Message) error {
			wg <- message
			return nil
		})

		err := c.LoadRoutes(ctx, strings.NewReader(simpleYAML))
		assert.Nil(t, err)

		select {
		case msg := <-wg:
			a, ok := msg.Annotation(timer.AnnotationTimerFiredCount)
			assert.True(t, ok)
			assert.Equal(t, "1", a)
			assert.Equal(t, content, msg.Content())
		case <-time.After(60 * time.Second):
			assert.Fail(t, "timeout")
		}
	})
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
	support.Run(t, "run", func(t *testing.T, ctx context.Context) {
		t.Helper()

		wg := make(chan camel.Message)

		c := camel.GetContext(ctx)

		c.Registry().Set("consumer-1", func(_ context.Context, message camel.Message) error {
			_ = message.SetSubject("consumer-1")
			return nil
		})
		c.Registry().Set("consumer-2", func(_ context.Context, message camel.Message) error {
			wg <- message
			return nil
		})

		err := c.LoadRoutes(ctx, strings.NewReader(simpleWASM))
		assert.Nil(t, err)

		select {
		case msg := <-wg:
			assert.Equal(t, "consumer-1", msg.GetSubject())

			c, ok := msg.Content().([]byte)
			assert.True(t, ok)
			assert.Equal(t, "hello from wasm", string(c))

		case <-time.After(5 * time.Second):
			assert.Fail(t, "timeout")
		}
	})
}

const simpleError = `
- route:
    from:
      uri: "timer:foo"
      steps:
        - process:
            ref: "panic"
`

func TestSimpleError(t *testing.T) {
	t.Skip("TODO")

	l, err := zap.NewDevelopment()
	assert.Nil(t, err)

	camelContext := core.NewContext(l, camelc.WithLogErrorHandler())
	ctx := context.WithValue(context.Background(), camel.ContextKeyCamelContext, camelContext)

	assert.NotNil(t, camelContext)

	defer func() {
		_ = camelContext.Close(ctx)
	}()

	wg := make(chan camel.Message)

	c := camel.GetContext(ctx)
	c.Registry().Set("panic", func(_ context.Context, message camel.Message) error {
		return errors.New("foo")
	})

	err = c.LoadRoutes(ctx, strings.NewReader(simpleError))
	assert.Nil(t, err)

	select {
	case msg := <-wg:
		a, ok := msg.Annotation(timer.AnnotationTimerFiredCount)
		assert.True(t, ok)
		assert.Equal(t, "1", a)
	case <-time.After(5 * time.Second):
		assert.Fail(t, "timeout")
	}
}
