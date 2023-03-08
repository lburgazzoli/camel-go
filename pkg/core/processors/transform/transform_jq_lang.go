package transform

import (
	"context"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
)

type LanguageJq struct {
	Expression string `yaml:"expression"`
}

func newJqProcessor(_ context.Context, definition *LanguageJq) (languageProcessor, error) {
	if definition.Expression == "" {
		return nil, camelerrors.MissingParameterf("jq.expression", "failure processing %s", TAG)
	}

	return &jqProcessor{}, nil
}

type jqProcessor struct {
}

func (p *jqProcessor) Process(ctx context.Context, m camel.Message) (camel.Message, error) {
	return nil, camelerrors.NotImplemented("todo")
}
