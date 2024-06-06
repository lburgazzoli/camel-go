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

func New(opts ...OptionFn) *Language {
	answer := &Language{}

	for _, o := range opts {
		o(answer)
	}

	return answer
}

type Language struct {
	Jq       *jq.Jq             `yaml:"jq,omitempty"`
	Mustache *mustache.Mustache `yaml:"mustache,omitempty"`
	Wasm     *wasm.Wasm         `yaml:"wasm,omitempty"`
	Constant *constant.Constant `yaml:"constant,omitempty"`
}

//nolint:dupl
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

//nolint:dupl
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

//nolint:dupl
func (l *Language) Transformer(ctx context.Context, camelContext camel.Context) (camel.Transformer, error) {
	switch {
	case l.Wasm != nil:
		p, err := l.Wasm.Transformer(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		return p, nil

	case l.Mustache != nil:
		p, err := l.Mustache.Transformer(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		return p, nil

	case l.Jq != nil:
		p, err := l.Jq.Transformer(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		return p, nil

	case l.Constant != nil:
		p, err := l.Constant.Transformer(ctx, camelContext)
		if err != nil {
			return nil, err
		}

		return p, nil
	default:
		return nil, camelerrors.MissingParameter("wasm || mustache || jq || constant", "missing language")
	}
}
