package transform

import (
	"bytes"
	"context"

	"github.com/cbroglie/mustache"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
)

type LanguageMustache struct {
	Template string `yaml:"template"`
}

func newMustacheProcessor(_ context.Context, definition *LanguageMustache) (languageProcessor, error) {
	if definition.Template == "" {
		return nil, camelerrors.MissingParameterf("mustache.template", "failure processing %s", TAG)
	}

	t, err := mustache.ParseString(definition.Template)
	if err != nil {
		return nil, err
	}

	return &mustacheProcessor{t: t}, nil
}

type mustacheProcessor struct {
	t *mustache.Template
}

func (p *mustacheProcessor) Process(ctx context.Context, m camel.Message) (camel.Message, error) {

	var buf bytes.Buffer

	err := p.t.FRender(&buf, map[string]interface{}{
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
		},
	})
	if err != nil {
		return nil, err
	}

	m.SetContent(buf.Bytes())

	return m, nil
}
