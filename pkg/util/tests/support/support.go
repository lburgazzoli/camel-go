package support

import (
	"bytes"
	"context"
	"github.com/lburgazzoli/camel-go/pkg/logger"
	"strings"
	"sync"
	"testing"
	"text/template"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core"
	"github.com/stretchr/testify/assert"
)

var once sync.Once

func Run(t *testing.T, name string, fn func(*testing.T, context.Context)) {
	t.Helper()

	once.Do(func() {
		logger.Init(logger.Options{
			Development: true,
		})
	})

	t.Run(name, func(t *testing.T) {
		camelContext := core.NewContext(logger.L)
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
