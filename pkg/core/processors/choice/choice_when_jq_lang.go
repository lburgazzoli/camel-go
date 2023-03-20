package choice

import (
	"context"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
)

type LanguageJq struct {
	Expression string `yaml:"expression"`
}

func newJqPredicate(_ context.Context, definition *LanguageJq) (languagePredicate, error) {
	if definition.Expression == "" {
		return nil, camelerrors.MissingParameterf("jq.expression", "failure processing %s", TAG)
	}

	return &jqPredicate{}, nil
}

type jqPredicate struct {
}

func (p *jqPredicate) Matches(ctx context.Context, m camel.Message) (bool, error) {
	return false, camelerrors.NotImplemented("todo")
}
