package wasm

import (
	"context"
	"errors"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
)

type Predicate struct {
	Function
}

func (p *Processor) Test(ctx context.Context, message camel.Message) (bool, error) {
	err := p.invoke(ctx, message)

	if err != nil {
		if errors.Is(err, ErrPredicateDoesNotMatch) {
			return false, nil
		}

		if errors.Is(err, ErrPredicateMatches) {
			return true, nil
		}

		return false, err
	}

	return false, nil
}
