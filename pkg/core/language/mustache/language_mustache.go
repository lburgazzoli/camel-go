package mustache

import (
	"bytes"
	"context"

	"github.com/cbroglie/mustache"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/pkg/errors"
)

type Mustache struct {
	Template string `yaml:"template"`
}

func (l *Mustache) Predicate(_ context.Context, _ camel.Context) (camel.Predicate, error) {
	return nil, camelerrors.NotImplemented("NotSupported")
}

func (l *Mustache) Processor(_ context.Context, _ camel.Context) (camel.Processor, error) {
	if l.Template == "" {
		return nil, camelerrors.MissingParameterf("mustache.template", "failure configuring jq processor")
	}

	t, err := mustache.ParseString(l.Template)
	if err != nil {
		return nil, err
	}

	p := func(ctxm context.Context, m camel.Message) error {
		answer, err := l.compute(t, m)
		if err != nil {
			return errors.Wrap(err, "error processing input")
		}

		m.SetContent(answer)

		return nil
	}

	return p, nil
}

func (l *Mustache) Transformer(_ context.Context, _ camel.Context) (camel.Transformer, error) {
	if l.Template == "" {
		return nil, camelerrors.MissingParameterf("mustache.template", "failure configuring jq processor")
	}

	t, err := mustache.ParseString(l.Template)
	if err != nil {
		return nil, err
	}

	p := func(ctxm context.Context, m camel.Message) (any, error) {
		return l.compute(t, m)
	}

	return p, nil
}

func (l *Mustache) compute(template *mustache.Template, m camel.Message) (any, error) {
	var buf bytes.Buffer

	err := template.FRender(&buf, map[string]interface{}{
		"message": map[string]interface{}{
			"id":             m.ID(),
			"source":         m.Source(),
			"type":           m.Type(),
			"subject":        m.Subject(),
			"content-type":   m.ContentType(),
			"content-schema": m.ContentSchema(),
			"time":           m.Time(),
			"content":        m.Content(),
			"attributes":     m.Attributes(),
			"header":         m.Headers(),
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "error processing input")
	}

	return buf.Bytes(), nil
}
