package wasm

import (
	"context"

	"github.com/knqyf263/go-plugin/types/known/timestamppb"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	pp "github.com/lburgazzoli/camel-go/pkg/wasm/plugin/processor"
)

type Processor struct {
	Function
}

func (p *Processor) Process(ctx context.Context, message camel.Message) error {

	camelContext := camel.ExtractContext(ctx)

	content := pp.Message{
		Id:            message.GetID(),
		Source:        message.GetSource(),
		Type:          message.GetType(),
		Subject:       message.GetSubject(),
		ContentType:   message.GetDataContentType(),
		ContentSchema: message.GetDataSchema(),
		Time:          timestamppb.New(message.GetTime()),
		Attributes:    make(map[string]string),
		Annotations:   make(map[string]string),
	}

	// TODO:fix annotation/attributes
	message.ForEachAnnotation(func(k string, v string) {
		content.Annotations[k] = v
	})

	_, err := camelContext.TypeConverter().Convert(message.Content(), &content.Data)
	if err != nil {
		return err
	}

	err = p.invoke(ctx, &content)
	if err != nil {
		return err
	}

	_ = message.SetID(content.Id)
	_ = message.SetSource(content.Source)
	_ = message.SetType(content.Type)
	_ = message.SetSubject(content.Subject)
	_ = message.SetDataContentType(content.ContentType)
	_ = message.SetDataSchema(content.ContentSchema)
	_ = message.SetTime(content.Time.AsTime())

	message.SetContent(content.Data)

	for k, v := range content.Annotations {
		message.SetAnnotation(k, v)
	}

	return nil
}
