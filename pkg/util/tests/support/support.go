package support

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"testing"
	"text/template"

	"github.com/onsi/gomega"

	"github.com/lburgazzoli/camel-go/pkg/logger"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core"
	"github.com/stretchr/testify/assert"
)

var once sync.Once

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

type Test interface {
	T() *testing.T
	Ctx() context.Context
	Camel() camel.Context

	gomega.Gomega
}

type T struct {
	*gomega.WithT

	t *testing.T

	//nolint:containedctx
	ctx context.Context
}

func With(t *testing.T) Test {
	t.Helper()

	once.Do(func() {
		logger.Init(logger.Options{
			Development: true,
		})
	})

	camelContext := core.NewContext(logger.L)
	assert.NotNil(t, camelContext)

	g := gomega.NewWithT(t)

	ctx := context.Background()
	ctx = context.WithValue(ctx, camel.ContextKeyCamelContext, camelContext)

	t.Cleanup(func() {
		err := camelContext.Close(ctx)
		g.Expect(err).NotTo(gomega.HaveOccurred())
	})

	if deadline, ok := t.Deadline(); ok {
		withDeadline, cancel := context.WithDeadline(ctx, deadline)
		t.Cleanup(cancel)

		ctx = withDeadline
	}

	return &T{
		WithT: g,
		t:     t,
		ctx:   ctx,
	}
}

func (t *T) T() *testing.T {
	return t.t
}

func (t *T) Ctx() context.Context {
	return t.ctx
}

func (t *T) Camel() camel.Context {
	return camel.ExtractContext(t.ctx)
}
