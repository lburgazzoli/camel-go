package support

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"text/template"

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

func LoadRoutes(ctx context.Context, route string, params any) error {

	c := camel.ExtractContext(ctx)

	tmpl, err := template.New("route").Parse(route)
	if err != nil {
		return err
	}

	var doc bytes.Buffer

	err = tmpl.Execute(&doc, params)
	if err != nil {
		return err
	}

	return c.LoadRoutes(ctx, strings.NewReader(doc.String()))
}
