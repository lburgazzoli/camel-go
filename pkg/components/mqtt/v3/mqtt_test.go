package v3

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// test support.
	"github.com/lburgazzoli/camel-go/pkg/util/tests/containers/mqtt"
	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"

	// enable components.
	_ "github.com/lburgazzoli/camel-go/pkg/components/log"
	_ "github.com/lburgazzoli/camel-go/pkg/components/timer"
	_ "github.com/lburgazzoli/camel-go/pkg/components/wasm"
	_ "github.com/twmb/franz-go/pkg/kgo"

	// enable processors.
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/process"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/to"
)

const simpleMQTT = `
- route:
    from:
      uri: "mqtt-v3:camel/iot"
      parameters:
        broker: "{{.broker}}"
      steps:
        - to:
            uri: "log:info"
        - process:
            ref: "consumer-1"
`

func TestSimpleMQTT(t *testing.T) {
	support.Run(t, "run", func(t *testing.T, ctx context.Context) {
		t.Helper()

		content := uuid.New()
		wg := make(chan camel.Message)

		conf, err := filepath.Abs("../../../../etc/support/mqtt/mosquitto.conf")
		require.NoError(t, err)
		require.FileExists(t, conf)

		container, err := mqtt.NewContainer(ctx, mqtt.WithConfig(conf))
		if err != nil {
			t.Error(err)
		}

		defer func() {
			if err := container.Stop(ctx); err != nil {
				t.Fatal(err.Error())
			}
		}()

		require.NoError(t, container.Start(ctx))

		cl, err := container.Client(ctx)
		require.NoError(t, err)

		c := camel.ExtractContext(ctx)
		c.Registry().Set("consumer-1", func(_ context.Context, message camel.Message) error {
			wg <- message
			return nil
		})

		broker, err := container.Broker(ctx)
		require.NoError(t, err)

		err = support.LoadRoutes(ctx, simpleMQTT, map[string]string{
			"broker": broker,
		})

		require.NoError(t, err)

		token := cl.Publish("camel/iot", 0, true, content)
		token.Wait()
		require.NoError(t, token.Error())

		select {
		case msg := <-wg:
			c, ok := msg.Content().([]byte)
			assert.True(t, ok)
			assert.Equal(t, content, string(c))

		case <-time.After(10 * time.Second):
			assert.Fail(t, "timeout")
		}
	})
}
