////go:build steps_process || steps_all

package process

import (
	"testing"
	"time"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"

	"github.com/stretchr/testify/assert"
)

func TestSimpleProcessor(t *testing.T) {
	content := uuid.New()
	wg := make(chan api.Message)

	c := core.NewContext()
	assert.NotNil(t, c)

	c.Registry().Set("p", func(message api.Message) {
		message.SetContent(content)
	})

	v := support.NewChannelVerticle(wg)
	err := c.Spawn(v)

	assert.Nil(t, err)

	p := Process{Identity: uuid.New(), Ref: "p"}
	p.Next(v.ID())

	id, err := p.Reify(c)
	assert.Nil(t, err)
	assert.NotNil(t, id)

	msg, err := message.New()
	assert.Nil(t, err)

	assert.Nil(t, c.Send(id, msg))

	select {
	case msg := <-wg:
		assert.Equal(t, content, msg.Content())
	case <-time.After(5 * time.Second):
		assert.Fail(t, "timeout")
	}
}
