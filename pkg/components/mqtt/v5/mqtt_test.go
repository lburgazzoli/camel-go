package v5

import (
	"context"
	"path/filepath"
	"strings"
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
      uri: "mqtt:camel/iot"
      parameters:
        broker: "${mqtt.broker}"
      steps:
        - to:
            uri: "log:info"
        - process:
            ref: "consumer-1"
`

func TestSimpleMQTT(t *testing.T) {
	g := support.With(t)

	content := uuid.New()
	wg := make(chan camel.Message)

	conf, err := filepath.Abs("../../../../etc/support/mqtt/mosquitto.conf")
	require.NoError(t, err)
	require.FileExists(t, conf)

	container, err := mqtt.NewContainer(g.Ctx(), mqtt.WithConfig(conf))
	if err != nil {
		t.Error(err)
	}

	defer func() {
		if err := container.Stop(g.Ctx()); err != nil {
			t.Fatal(err.Error())
		}
	}()

	require.NoError(t, container.Start(g.Ctx()))

	cl, err := container.Client(g.Ctx())
	require.NoError(t, err)

	broker, err := container.Broker(g.Ctx())
	require.NoError(t, err)

	c := camel.ExtractContext(g.Ctx())

	err = c.Properties().Add(map[string]any{"mqtt.broker": broker})
	require.NoError(t, err)

	c.Registry().Set("consumer-1", func(_ context.Context, message camel.Message) error {
		wg <- message
		return nil
	})

	err = c.LoadRoutes(g.Ctx(), strings.NewReader(simpleMQTT))
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
}
