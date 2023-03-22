package jq

import (
	"context"
	"fmt"
	"strconv"

	"github.com/itchyny/gojq"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/pkg/errors"
)

const (
	AnnotationJqResults = "camel.apache.org/jq.results"
)

type Jq struct {
	Expression string `yaml:"expression"`
}

func (l *Jq) run(
	ctx context.Context,
	camelContext camel.Context,
	query *gojq.Query,
	m camel.Message,
) (gojq.Iter, error) {

	var input camel.RawJSON

	ok, err := camelContext.TypeConverter().Convert(m.Content(), &input)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("unable to convert %+v to camel.RawJSON", m.Content())
	}

	return query.RunWithContext(ctx, map[string]interface{}(input)), nil
}

func (l *Jq) Predicate(ctx context.Context, camelContext camel.Context) (camel.Predicate, error) {
	if l.Expression == "" {
		return nil, camelerrors.MissingParameterf("jq.expression", "failure configuring jq predicate")
	}

	query, err := gojq.Parse(l.Expression)
	if err != nil {
		return nil, err
	}

	p := func(ctxm context.Context, m camel.Message) (bool, error) {
		it, err := l.run(ctx, camelContext, query, m)
		if err != nil {
			return false, err
		}

		v, ok := it.Next()
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

	return p, nil
}

func (l *Jq) Processor(ctx context.Context, camelContext camel.Context) (camel.Processor, error) {
	if l.Expression == "" {
		return nil, camelerrors.MissingParameterf("jq.expression", "failure configuring jq processor")
	}

	query, err := gojq.Parse(l.Expression)
	if err != nil {
		return nil, err
	}

	p := func(ctxm context.Context, m camel.Message) error {
		it, err := l.run(ctx, camelContext, query, m)
		if err != nil {
			return err
		}

		out := make([]interface{}, 0, 1)

		for {
			v, ok := it.Next()
			if !ok {
				break
			}

			if err, ok := v.(error); ok {
				return errors.Wrap(err, "error processing input")
			}

			out = append(out, v)
		}

		m.SetAnnotation(AnnotationJqResults, strconv.Itoa(len(out)))

		if len(out) == 1 {
			m.SetContent(out[0])
		} else if len(out) > 1 {
			m.SetContent(out)
		}

		return nil
	}

	return p, nil
}
