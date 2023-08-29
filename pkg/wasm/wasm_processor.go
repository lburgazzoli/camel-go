package wasm

import (
	"context"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	pp "github.com/lburgazzoli/camel-go/pkg/wasm/plugin/processor"
)

type Processor struct {
	Function
}

func (p *Processor) Process(ctx context.Context, message camel.Message) error {

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
		return err
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
		return err
	}

	_, err := camelContext.TypeConverter().Convert(message.Content(), &content.Data)
	if err != nil {
		return err
	}

	err = p.invoke(ctx, &content, &content)
	if err != nil {
		return err
	}

	message.SetSource(content.Source)
	message.SetType(content.Type)
	message.SetSubject(content.Subject)
	message.SetContentType(content.ContentType)
	message.SetContentSchema(content.ContentSchema)
	message.SetContent(content.Data)

	for k, v := range content.Headers {
		message.SetHeader(k, v)
	}
	for k, v := range content.Attributes {
		message.SetAttribute(k, v)
	}

	return nil
}
