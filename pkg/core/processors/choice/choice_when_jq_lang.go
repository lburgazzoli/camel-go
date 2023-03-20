package choice

import (
	"context"
	"fmt"

	"github.com/itchyny/gojq"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
)

type LanguageJq struct {
	Expression string `yaml:"expression"`
}

func newJqPredicate(_ context.Context, camelContext camel.Context, definition *LanguageJq) (languagePredicate, error) {
	if definition.Expression == "" {
		return nil, camelerrors.MissingParameterf("jq.expression", "failure processing %s", TAG)
	}

	query, err := gojq.Parse(definition.Expression)
	if err != nil {
		return nil, err
	}

	return &jqPredicate{
		camelContext: camelContext,
		query:        query,
	}, nil
}

type jqPredicate struct {
	camelContext camel.Context
	query        *gojq.Query
}

func (p *jqPredicate) Matches(ctx context.Context, m camel.Message) (bool, error) {
	var input camel.RawJSON

	ok, err := p.camelContext.TypeConverter().Convert(m.Content(), &input)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, fmt.Errorf("unable to convert %+v to amel.RawJSON", m.Content())
	}

	in := map[string]interface{}(input)
	iter := p.query.RunWithContext(ctx, in)
	v, ok := iter.Next()
	if !ok {
		return false, nil
	}

	if err, ok := v.(error); ok {
		return false, err
	}

	if match, ok := v.(bool); ok {
		return match, nil
	}

	return false, nil
}
