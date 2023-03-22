package kafka

import (
	"context"
	"crypto/tls"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/twmb/franz-go/plugin/kzap"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/pkg/errors"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
)

type Producer struct {
	api.WithOutputs

	id       string
	endpoint *Endpoint
	client   *kgo.Client
	tc       api.TypeConverter
}

func (p *Producer) ID() string {
	return p.id
}

func (p *Producer) Endpoint() api.Endpoint {
	return p.endpoint
}

func (p *Producer) Start(context.Context) error {

	opts := make([]kgo.Opt, 0)
	opts = append(opts, kgo.SeedBrokers(strings.Split(p.endpoint.config.Brokers, ",")...))

	if p.endpoint.config.User != "" && p.endpoint.config.Password != "" {
		tlsDialer := &tls.Dialer{NetDialer: &net.Dialer{Timeout: 10 * time.Second}}
		authMechanism := plain.Auth{User: p.endpoint.config.User, Pass: p.endpoint.config.Password}.AsMechanism()

		opts = append(opts, kgo.SASL(authMechanism))
		opts = append(opts, kgo.Dialer(tlsDialer.DialContext))
		opts = append(opts, kgo.WithLogger(kzap.New(p.endpoint.Logger())))
	}

	cl, err := kgo.NewClient(opts...)

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

		for _, o := range p.Outputs() {
			ac.Send(o, msg)
		}
	}
}

func (p *Producer) publish(ctx context.Context, msg api.Message) {
	record := &kgo.Record{}
	record.Topic = p.endpoint.config.Remaining
	record.Headers = []kgo.RecordHeader{
		{Key: "event.id", Value: []byte(msg.GetID())},
		{Key: "event.type", Value: []byte(msg.GetType())},
	}

	if s := msg.GetSubject(); s != "" {
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
		msg.SetAnnotation(AnnotationOffset, strconv.FormatInt(r.Offset, 10))
		msg.SetAnnotation(AnnotationPartition, strconv.FormatInt(int64(r.Partition), 10))
	}
}
