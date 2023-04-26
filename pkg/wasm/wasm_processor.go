package wasm

import (
	"context"

	"github.com/knqyf263/go-plugin/types/known/timestamppb"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelmsg "github.com/lburgazzoli/camel-go/pkg/core/message"
	pp "github.com/lburgazzoli/camel-go/pkg/wasm/plugin/processor"
)

type Function struct {
	processor pp.Processors
}

// Invoke invoke a function.
func (f *Function) Invoke(ctx context.Context, message camel.Message) (camel.Message, error) {

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
		return message, err
	}

	result, err := f.processor.Process(ctx, &content)
	if err != nil {
		return nil, err
	}

	msg, err := camelmsg.New()
	if err != nil {
		return nil, err
	}

	_ = msg.SetID(result.Id)
	_ = msg.SetSource(result.Source)
	_ = msg.SetType(result.Type)
	_ = msg.SetSubject(result.Subject)
	_ = msg.SetDataContentType(result.ContentType)
	_ = msg.SetDataSchema(result.ContentSchema)
	_ = msg.SetTime(result.Time.AsTime())

	msg.SetContent(result.Data)

	for k, v := range result.Annotations {
		msg.SetAnnotation(k, v)
	}

	return msg, nil
}
