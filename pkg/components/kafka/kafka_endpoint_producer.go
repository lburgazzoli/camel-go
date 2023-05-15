////go:build components_kafka || components_all

package kafka

import (
	"context"
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

func (p *Producer) publish(ctx context.Context, msg api.Message) {
	record := &kgo.Record{}
	record.Topic = p.endpoint.config.Remaining
	record.Headers = []kgo.RecordHeader{
		{Key: "event.id", Value: []byte(msg.ID())},
		{Key: "event.type", Value: []byte(msg.Type())},
	}

	if s := msg.Subject(); s != "" {
		record.Key = []byte(s)
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
