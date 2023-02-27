package kafka

import (
	"context"
	"strings"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/pkg/errors"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	api.WithOutputs

	id       string
	endpoint *Endpoint
	client   *kgo.Client
}

func (p *Producer) ID() string {
	return p.id
}

func (p *Producer) Endpoint() api.Endpoint {
	return p.endpoint
}

func (p *Producer) Start() error {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(strings.Split(p.endpoint.config.Brokers, ",")...),
	)

	if err != nil {
		return err
	}

	p.client = cl

	return nil
}

func (p *Producer) Stop() error {
	if p.client != nil {
		p.client.Close()
		p.client = nil
	}

	return nil
}

func (p *Producer) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		_ = p.Start()
	case *actor.Stopping:
		_ = p.Stop()
	case api.Message:
		if err := p.publish(msg); err != nil {
			panic(err)
		}
	}
}

func (p *Producer) publish(msg api.Message) error {
	record := &kgo.Record{}
	record.Topic = p.endpoint.config.Topics
	record.Headers = []kgo.RecordHeader{
		{Key: "event.id", Value: []byte(msg.GetID())},
		{Key: "event.type", Value: []byte(msg.GetType())},
	}

	// TODO: must implement type converters
	switch v := msg.Content().(type) {
	case []byte:
		record.Value = v
	case string:
		record.Value = []byte(v)
	default:
		panic("unsupported content type")
	}

	// TODO: must get a context.Context
	if err := p.client.ProduceSync(context.TODO(), record).FirstErr(); err != nil {
		return errors.Wrap(err, "record had a produce error while synchronously producing")
	}

	return nil
}
