// //go:build steps_process || steps_all

package setbody

import (
	"context"
	"testing"
	"time"

	"github.com/lburgazzoli/camel-go/pkg/core/language"
	"github.com/lburgazzoli/camel-go/pkg/core/language/constant"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/stretchr/testify/require"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
)

func TestSetBody(t *testing.T) {
	g := support.With(t)
	c := camel.ExtractContext(g.Ctx())

	content := uuid.New()

	c.Registry().Set("p", func(_ context.Context, message camel.Message) error {
		message.SetContent(content)
		return nil
	})

	p := New()
	p.Language = language.Language{
		Constant: &constant.Constant{
			Value: content,
		},
	}

	pv, err := p.Reify(g.Ctx())
	require.NoError(t, err)
	require.NotNil(t, pv)

	pvp, err := c.Spawn(pv)
	require.NoError(t, err)
	require.NotNil(t, pvp)

	msg := c.NewMessage()

	res, err := c.RequestTo(pvp, msg, 1*time.Second)
	require.NoError(t, err)
	require.Equal(t, content, res.Content())
}
