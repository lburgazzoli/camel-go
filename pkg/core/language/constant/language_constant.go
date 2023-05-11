package constant

import (
	"context"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
)

type Constant struct {
	Value string `yaml:"value"`
}

func (l *Constant) Predicate(_ context.Context, _ camel.Context) (camel.Predicate, error) {
	if l.Value == "" {
		return nil, camelerrors.MissingParameterf("constant.value", "failure configuring constant predicate")
	}

	p := func(ctx context.Context, message camel.Message) (bool, error) {
		return l.Value == "true", nil
	}

	return p, nil
}

func (l *Constant) Processor(_ context.Context, _ camel.Context) (camel.Processor, error) {
	if l.Value == "" {
		return nil, camelerrors.MissingParameterf("constant.value", "failure configuring constant processor")
	}

	p := func(ctxm context.Context, m camel.Message) error {
		m.SetContent(l.Value)

		return nil
	}

	return p, nil
}
