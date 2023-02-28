////go:build steps_transform || steps_all

package transform

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSimpleTransform(t *testing.T) {
	wg := make(chan api.Message)

	c := core.NewContext()
	assert.NotNil(t, c)

	v := support.NewChannelVerticle(wg)
	err := c.Spawn(v)

	assert.Nil(t, err)

	p := Transform{
		Identity: uuid.New(),
		Language: Language{
			Wasm: &LanguageWasm{
				Path: "../../../../etc/fn/simple_process.wasm",
			},
		}}

	p.Next(v.ID())

	id, err := p.Reify(c)
	assert.Nil(t, err)
	assert.NotNil(t, id)

	msg, err := message.New()
	assert.Nil(t, err)

	assert.Nil(t, c.Send(id, msg))

	select {
	case msg := <-wg:
		c, ok := msg.Content().([]byte)
		assert.True(t, ok)
		assert.Equal(t, "hello from wasm", string(c))
	case <-time.After(5 * time.Second):
		assert.Fail(t, "timeout")
	}
}
