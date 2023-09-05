// //go:build components_dapr_pubsub || components_all

package pubsub

import (
	"context"
	"log/slog"
	"strings"

	"github.com/lburgazzoli/camel-go/pkg/components"

	"github.com/asynkron/protoactor-go/actor"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/pkg/errors"
)

type Producer struct {
	components.DefaultProducer

	endpoint *Endpoint
	client   dapr.Client
	tc       api.TypeConverter

	pubsubName string
	topicName  string
}

func (p *Producer) Endpoint() api.Endpoint {
	return p.endpoint
}

func (p *Producer) Start(context.Context) error {
	ct := strings.Split(p.endpoint.config.Remaining, "/")
	if len(ct) != 2 {
		return camelerrors.MissingParameter("pubsubName/topicName", "missing pubsubName/topicName")
	}

	p.pubsubName = ct[0]
	p.topicName = ct[1]

	cl, err := dapr.NewClient()
	if err != nil {
		return err
	}

	p.client = cl

	p.Logger().Info("started",
		slog.Group(
			"target",
			slog.String("pubsubName", p.pubsubName),
			slog.String("topic", p.topicName),
		),
	)

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

		if msg.Error() != nil {
			panic(msg.Error())
		}

		ac.Request(ac.Parent(), msg)
	}
}

func (p *Producer) publish(ctx context.Context, msg api.Message) {
	if p.client == nil {
		msg.SetError(errors.New("client is not yet ready"))
		return
	}

	data := make([]byte, 0)

	if _, err := p.tc.Convert(msg.Content(), &data); err != nil {
		msg.SetError(errors.Wrap(err, "error converting content to []byte"))
		return
	}

	opts := make([]dapr.PublishEventOption, 0)
	if msg.ContentType() != "" {
		opts = append(opts, dapr.PublishEventWithContentType(msg.ContentType()))
	}

	if err := p.client.PublishEvent(ctx, p.pubsubName, p.topicName, data, opts...); err != nil {
		msg.SetError(err)
		return
	}
}
