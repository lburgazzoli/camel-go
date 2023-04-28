package support

import (
	"context"
	"testing"

	"go.uber.org/zap"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Run(t *testing.T, name string, fn func(*testing.T, context.Context)) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		l, err := zap.NewDevelopment()
		assert.Nil(t, err)

		camelContext := core.NewContext(l)
		ctx := context.WithValue(context.Background(), camel.ContextKeyCamelContext, camelContext)

		assert.NotNil(t, camelContext)

		defer func() {
			_ = camelContext.Close(ctx)
		}()

		fn(t, ctx)
	})
}
