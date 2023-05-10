package constant

import (
	"context"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
)

type Constant struct {
	Value string `yaml:"value"`
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
