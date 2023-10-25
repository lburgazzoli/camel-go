package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/xid"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// enable components.
	_ "github.com/lburgazzoli/camel-go/pkg/components/timer"

	// enable processors.
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/process"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/to"
)

const simpleHTTPGet = `
- route:
    from:
      uri: "timer:foo"
      steps:
        - process:
            ref: "consumer-1"
        - to:
            uri: "{{.URL}}"
        - process:
            ref: "consumer-2"
`

func TestSimpleHTTPGet(t *testing.T) {
	g := support.With(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		answer := map[string]any{
			"uuid": xid.New().String(),
		}

		data, err := json.Marshal(answer)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		_, err = io.WriteString(w, string(data))
		require.NoError(t, err)
	}))

	defer func() {
		ts.Close()
	}()

	wg := make(chan camel.Message)

	c := camel.ExtractContext(g.Ctx())

	c.Registry().Set("consumer-1", func(_ context.Context, message camel.Message) error {
		message.SetHeader("Accept", "application/json")
		return nil
	})
	c.Registry().Set("consumer-2", func(_ context.Context, message camel.Message) error {
		wg <- message
		return nil
	})

	err := support.LoadRoutes(
		g.Ctx(),
		simpleHTTPGet,
		map[string]string{
			"URL": ts.URL,
		},
	)

	require.NoError(t, err)

	select {
	case msg := <-wg:
		c, ok := msg.Content().([]byte)
		require.True(t, ok)
		require.NotEmpty(t, c)

		ct, ok := msg.Header("Content-Type")
		require.True(t, ok)
		require.Equal(t, "application/json", ct)

		sc, ok := msg.Attribute(AttributeStatusCode)
		require.True(t, ok)
		require.Equal(t, 200, sc)

		data, err := message.ContentAsBytes(msg)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		m := make(map[string]any)
		err = json.Unmarshal(data, &m)
		require.NoError(t, err)
		require.NotEmpty(t, m["uuid"])

	case <-time.After(60 * time.Second):
		assert.Fail(t, "timeout")
	}
}
