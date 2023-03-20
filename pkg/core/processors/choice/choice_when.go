package choice

import (
	"context"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

type When struct {
	processors.DefaultVerticle `yaml:",inline"`
	Language                   `yaml:",inline"`

	predicate languagePredicate
	Steps     []processors.Step `yaml:"steps,omitempty"`
}

func (w *When) Configure(ctx context.Context, camelContext camel.Context) error {
	w.DefaultVerticle.SetContext(camelContext)

	switch {
	case w.Jq != nil:
		p, err := newJqPredicate(ctx, camelContext, w.Jq)
		if err != nil {
			return err
		}

		w.predicate = p
	default:
		return camelerrors.MissingParameterf("jq", "failure processing %s", TAG)
	}

	return nil
}

func (w *When) Matches(ctx context.Context, msg camel.Message) (bool, error) {
	if w.predicate == nil {
		return false, camelerrors.InternalErrorf("not configured")
	}

	return w.predicate.Matches(ctx, msg)
}

type languagePredicate interface {
	Matches(context.Context, camel.Message) (bool, error)
}

type Language struct {
	Jq *LanguageJq `yaml:"jq,omitempty"`
}
