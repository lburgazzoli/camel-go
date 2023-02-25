////go:build components_all || components_timer || steps_all || steps_process

package engine

import (
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	cameltest "github.com/lburgazzoli/camel-go/test/support"
	"strings"

	"testing"
	"time"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/endpoint"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/from"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/process"
	"github.com/stretchr/testify/assert"

	_ "github.com/lburgazzoli/camel-go/pkg/components/timer"
)

func TestSimple(t *testing.T) {
	wg := make(chan api.Message)

	c := core.NewContext()
	assert.NotNil(t, c)

	c.Registry().Set("consumer", func(message api.Message) {
		wg <- message
	})

	f := from.From{
		Endpoint: endpoint.Endpoint{
			URI: "timer:foo",
			Parameters: map[string]interface{}{
				"interval": 1 * time.Second,
			},
		},
	}

	f.Next(cameltest.Reify(t, c, &process.Process{Ref: "consumer"}))

	fromPid, err := f.Reify(c)
	assert.Nil(t, err)
	assert.NotNil(t, fromPid)

	select {
	case msg := <-wg:
		a, ok := msg.Annotation("counter")
		assert.True(t, ok)
		assert.Equal(t, uint64(1), a)
	case <-time.After(5 * time.Second):
		assert.Fail(t, "timeout")
	}
}

const simpleRoute = `
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
	content := uuid.New()
	wg := make(chan api.Message)

	c := core.NewContext()
	assert.NotNil(t, c)

	c.Registry().Set("consumer-1", func(message api.Message) {
		message.SetContent(content)
	})
	c.Registry().Set("consumer-2", func(message api.Message) {
		wg <- message
	})

	err := c.LoadRoutes(strings.NewReader(simpleRoute))
	assert.Nil(t, err)

	select {
	case msg := <-wg:
		a, ok := msg.Annotation("counter")
		assert.True(t, ok)
		assert.Equal(t, uint64(1), a)
		//assert.Equal(t, content, msg.Content())
	case <-time.After(5 * time.Second):
		assert.Fail(t, "timeout")
	}
}
