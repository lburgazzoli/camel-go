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
		Id:            message.ID(),
		Source:        message.Source(),
		Type:          message.Type(),
		Subject:       message.Subject(),
		ContentType:   message.ContentType(),
		ContentSchema: message.ContentSchema(),
		Time:          timestamppb.New(message.Time()),
		Attributes:    make(map[string]string),
		Annotations:   make(map[string]string),
	}

	// TODO:fix annotation/attributes
	// message.EachAttribute(func(k string, v any) {
	//  content.Annotations[k] = v
	// })

	_, err := camelContext.TypeConverter().Convert(message.Content(), &content.Data)
	if err != nil {
		return err
	}

	err = p.invoke(ctx, &content)
	if err != nil {
		return err
	}

	message.SetSource(content.Source)
	message.SetType(content.Type)
	message.SetSubject(content.Subject)
	message.SetContentType(content.ContentType)
	message.SetContentSchema(content.ContentSchema)
	message.SetContent(content.Data)

	// TODO:fix annotation/attributes
	// for k, v := range content.Annotations {
	// 	 if err := message.SetAttribute(k, v); err != nil {
	// 	    	panic(err)
	//  	}
	// }

	return nil
}
