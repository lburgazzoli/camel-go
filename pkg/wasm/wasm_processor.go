package wasm

import (
	"context"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
)

type Processor struct {
	Function
}

func (p *Processor) Process(ctx context.Context, message camel.Message) error {
	err := p.invoke(ctx, message)
	if err != nil {
		return err
	}

	return nil
}
