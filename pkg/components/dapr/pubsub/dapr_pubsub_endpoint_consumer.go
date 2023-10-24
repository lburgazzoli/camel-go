// //go:build components_dapr_pubsub || components_all

package pubsub

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/lburgazzoli/camel-go/pkg/components/dapr"

	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dapr/go-sdk/service/common"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
)

type Consumer struct {
	components.DefaultConsumer

	endpoint *Endpoint
	running  atomic.Bool

	pubsubName string
	topicName  string
}

func (c *Consumer) Endpoint() camel.Endpoint {
	return c.endpoint
}

func (c *Consumer) Start(_ context.Context) error {
	if c.running.CompareAndSwap(false, true) {
		ct := strings.Split(c.endpoint.config.Remaining, "/")
		if len(ct) != 2 {
			return camelerrors.MissingParameter("pubsubName/topicName", "missing pubsubName/topicName")
		}

		c.pubsubName = ct[0]
		c.topicName = ct[1]

		sub := common.Subscription{
			PubsubName: c.pubsubName,
			Topic:      c.topicName,
			Route:      "/" + c.endpoint.ID() + "/" + c.ID(),
			Metadata: map[string]string{
				dapr.MetaRawPayload: strconv.FormatBool(c.endpoint.config.Raw),
			},
		}

		c.Logger().Debug("subscribing", slog.Group(
			"subscription",
			slog.String("pubsubName", sub.PubsubName),
			slog.String("topic", sub.Topic),
			slog.String("route", sub.Route),
			slog.Any("meta", sub.Metadata),
		))

		err := dapr.AddTopicEventHandler(&sub, c.handler)
		if err != nil {
			return err
		}

		c.Logger().Debug("subscribed", slog.Group(
			"subscription",
			slog.String("pubsubName", sub.PubsubName),
			slog.String("topic", sub.Topic),
			slog.String("route", sub.Route),
		))

		err = dapr.Start()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Consumer) Stop(context.Context) error {
	if c.running.CompareAndSwap(true, false) {
		return dapr.Stop()
	}

	return nil
}

func (c *Consumer) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		_ = c.Start(context.Background())
	case *actor.Stopping:
		_ = c.Stop(context.Background())
	case camel.Message:
		// ignore message,
		// TODO: may be used for ack
		break
	}
}

func (c *Consumer) handler(_ context.Context, e *common.TopicEvent) (bool, error) {
	component := c.endpoint.Component()
	camelCtx := component.Context()

	c.Logger().Debug("event received", slog.Group(
		"event",
		slog.String("pubsubName", e.PubsubName),
		slog.String("topic", e.Topic),
		slog.String("id", e.ID),
		slog.String("content-type", e.DataContentType),
	))

	m := camelCtx.NewMessage()

	m.SetType(e.Type)
	m.SetSource(e.Source)
	m.SetSubject(e.Subject)
	m.SetContentType(e.DataContentType)
	m.SetContent(e.RawData)

	m.SetAttribute(AttributeEventID, e.ID)
	m.SetAttribute(AttributePubSubName, e.PubsubName)
	m.SetAttribute(AttributePubSubTopic, e.Topic)

	if err := camelCtx.SendTo(c.Target(), m); err != nil {
		panic(err)
	}

	return false, nil
}
