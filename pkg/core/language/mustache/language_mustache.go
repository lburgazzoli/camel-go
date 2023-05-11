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
	return nil, camelerrors.NotImplemented("TODO")
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
		var buf bytes.Buffer

		err := t.FRender(&buf, map[string]interface{}{
			"message": map[string]interface{}{
				"id":                m.GetID(),
				"source":            m.GetSource(),
				"type":              m.GetType(),
				"subject":           m.GetSubject(),
				"data-content-type": m.GetDataContentType(),
				"data-schema":       m.GetDataSchema(),
				"time":              m.GetTime(),
				"content":           m.Content(),
				"annotations":       m.Annotations(),
				"extensions":        m.GetExtensions(),
			},
		})

		if err != nil {
			return errors.Wrap(err, "error processing input")
		}

		m.SetContent(buf.Bytes())

		return nil
	}

	return p, nil
}
