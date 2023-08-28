////go:build components_kafka || components_all

package kafka

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/pkg/errors"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	processors.DefaultVerticle

	endpoint *Endpoint
	client   *kgo.Client
	tc       api.TypeConverter
}

func (p *Producer) Endpoint() api.Endpoint {
	return p.endpoint
}

func (p *Producer) Start(context.Context) error {
	cl, err := p.endpoint.newClient()
	if err != nil {
		return err
	}

	p.client = cl

	return nil
}

func (p *Producer) Stop(context.Context) error {
	if p.client != nil {
		p.client.Close()
		p.client = nil
	}

	return nil
}

func (p *Producer) Receive(ac actor.Context) {
	switch msg := ac.Message().(type) {
	case *actor.Started:
		_ = p.Start(context.Background())
	case *actor.Stopping:
		_ = p.Stop(context.Background())
	case api.Message:
		p.publish(context.Background(), msg)

		// TODO: handle
		if msg.Error() != nil {
			panic(msg.Error())
		}

		ac.Request(ac.Parent(), msg)
	}
}

// publish produces a record that conforms to the CloudEvents binary-content-mode 1.0 spec.
func (p *Producer) publish(ctx context.Context, msg api.Message) {
	record := &kgo.Record{}
	record.Topic = p.endpoint.config.Remaining
	record.Headers = make([]kgo.RecordHeader, 5)

	record.Headers = append(record.Headers, kgo.RecordHeader{
		Key:   "ce_specversion",
		Value: []byte("1.0"),
	})

	// copy relevant attributes as ce headers
	_ = msg.EachAttribute(func(k string, v any) error {
		switch k {
		case api.MessageAttributeID:
			k = "id"
		case api.MessageAttributeTime:
			k = "time"
		case api.MessageAttributeSource:
			k = "source"
		case api.MessageAttributeContentType:
			k = "content-type"
		case api.MessageAttributeContentSchema:
			k = "datacontentschema"
		default:
			return nil
		}

		h, err := p.setHeader(k, v)
		if err != nil {
			msg.SetError(err)
			return nil
		}

		record.Headers = append(record.Headers, h)

		return nil

	})

	// copy remaining headers a standard headers
	_ = msg.EachHeader(func(k string, v any) error {
		h, err := p.setHeader(k, v)
		if err != nil {
			msg.SetError(err)
			return nil
		}

		record.Headers = append(record.Headers, h)

		return nil
	})

	if v := msg.Subject(); v != "" {
		record.Key = []byte(v)
	}

	_, err := p.tc.Convert(msg.Content(), &record.Value)
	if err != nil {
		msg.SetError(errors.Wrap(err, "error converting content to []byte"))
		return
	}

	result := p.client.ProduceSync(ctx, record)

	r, err := result.First()
	if err != nil {
		msg.SetError(errors.Wrap(err, "record had a produce error while synchronously producing"))
	}
	if r != nil {
		_ = msg.SetAttribute(AttributeOffset, strconv.FormatInt(r.Offset, 10))
		_ = msg.SetAttribute(AttributePartition, strconv.FormatInt(int64(r.Partition), 10))
	}
}

func (p *Producer) setHeader(k string, v any) (kgo.RecordHeader, error) {
	h := kgo.RecordHeader{}
	h.Key = k

	ok, err := p.tc.Convert(v, &h.Value)
	if err != nil {
		return h, err
	}
	if !ok {
		return h, fmt.Errorf("unable to convert value for header %s", k)
	}

	return h, nil
}
