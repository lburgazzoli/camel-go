package language

import (
	"context"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/core/language/constant"
	"github.com/lburgazzoli/camel-go/pkg/core/language/jq"
	"github.com/lburgazzoli/camel-go/pkg/core/language/mustache"
	"github.com/lburgazzoli/camel-go/pkg/core/language/wasm"
)

type Language struct {
	Jq       *jq.Jq             `yaml:"jq,omitempty"`
	Mustache *mustache.Mustache `yaml:"mustache,omitempty"`
	Wasm     *wasm.Wasm         `yaml:"wasm,omitempty"`
	Constant *constant.Constant `yaml:"constant,omitempty"`
}

func (l *Language) Processor(ctx context.Context, camelContext camel.Context) (camel.Processor, error) {

	switch {
	case l.Wasm != nil:
		p, err := l.Wasm.Processor(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		return p, nil

	case l.Mustache != nil:
		p, err := l.Mustache.Processor(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		return p, nil

	case l.Jq != nil:
		p, err := l.Jq.Processor(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		return p, nil

	case l.Constant != nil:
		p, err := l.Constant.Processor(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		return p, nil
	default:
		return nil, camelerrors.MissingParameter("wasm || mustache || jq || constant", "missing language")
	}
}

func (l *Language) Predicate(ctx context.Context, camelContext camel.Context) (camel.Predicate, error) {

	switch {
	case l.Wasm != nil:
		p, err := l.Wasm.Predicate(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		return p, nil

	case l.Mustache != nil:
		p, err := l.Mustache.Predicate(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		return p, nil

	case l.Jq != nil:
		p, err := l.Jq.Predicate(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		return p, nil

	case l.Constant != nil:
		p, err := l.Constant.Predicate(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		return p, nil
	default:
		return nil, camelerrors.MissingParameter("wasm || mustache || jq || constant", "missing language")
	}
}
