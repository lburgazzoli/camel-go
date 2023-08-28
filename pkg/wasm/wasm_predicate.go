package wasm

import (
	"context"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	pp "github.com/lburgazzoli/camel-go/pkg/wasm/plugin/processor"
)

type Predicate struct {
	Function
}

func (p *Processor) Test(ctx context.Context, message camel.Message) (bool, error) {

	camelContext := camel.ExtractContext(ctx)

	content := pp.Message{
		ID:            message.ID(),
		Source:        message.Source(),
		Type:          message.Type(),
		Subject:       message.Subject(),
		ContentType:   message.ContentType(),
		ContentSchema: message.ContentSchema(),
		Time:          message.Time(),
		Headers:       make(map[string][]byte),
		Attributes:    make(map[string][]byte),
	}

	if err := message.EachHeader(func(k string, v any) error {
		val := make([]byte, 0)

		_, err := camelContext.TypeConverter().Convert(v, &val)
		if err != nil {
			return err
		}

		content.Headers[k] = val

		return nil
	}); err != nil {
		return false, err
	}

	if err := message.EachAttribute(func(k string, v any) error {
		val := make([]byte, 0)

		_, err := camelContext.TypeConverter().Convert(v, &val)
		if err != nil {
			return err
		}

		content.Attributes[k] = val

		return nil
	}); err != nil {
		return false, err
	}

	_, err := camelContext.TypeConverter().Convert(message.Content(), &content.Data)
	if err != nil {
		return false, err
	}

	eval := pp.Evaluation{}

	err = p.invoke(ctx, &content, &eval)
	if err != nil {
		return false, err
	}

	return eval.Result, nil
}
