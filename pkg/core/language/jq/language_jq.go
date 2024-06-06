package jq

import (
	"context"
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"

	"github.com/itchyny/gojq"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/pkg/errors"
)

const (
	AnnotationJqResults = "camel.apache.org/jq.results"
)

type Definition struct {
	Expression string `yaml:"expression"`
}

type Jq struct {
	Definition `yaml:",inline"`
}

func (l *Jq) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		return l.UnmarshalText([]byte(value.Value))
	case yaml.MappingNode:
		return value.Decode(&l.Definition)
	default:
		return fmt.Errorf("unsupported node kind: %v (line: %d, column: %d)", value.Kind, value.Line, value.Column)
	}
}

func (l *Jq) UnmarshalText(text []byte) error {
	l.Expression = string(text)
	return nil
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
		answer, err := l.compute(ctx, camelContext, query, m)
		if err != nil {
			return err
		}

		m.SetContent(answer)

		return nil
	}

	return p, nil
}

func (l *Jq) Transformer(ctx context.Context, camelContext camel.Context) (camel.Transformer, error) {
	if l.Expression == "" {
		return nil, camelerrors.MissingParameterf("jq.expression", "failure configuring jq processor")
	}

	query, err := gojq.Parse(l.Expression)
	if err != nil {
		return nil, err
	}

	p := func(ctxm context.Context, m camel.Message) (any, error) {
		return l.compute(ctx, camelContext, query, m)
	}

	return p, nil
}

func (l *Jq) compute(ctx context.Context, camelContext camel.Context, query *gojq.Query, m camel.Message) (any, error) {
	it, err := l.run(ctx, camelContext, query, m)
	if err != nil {
		return nil, err
	}

	out := make([]interface{}, 0, 1)

	for {
		v, ok := it.Next()
		if !ok {
			break
		}

		if err, ok := v.(error); ok {
			return nil, errors.Wrap(err, "error processing input")
		}

		out = append(out, v)
	}

	m.SetAttribute(AnnotationJqResults, strconv.Itoa(len(out)))

	if len(out) == 1 {
		return out[0], nil
	}

	return out, nil
}
