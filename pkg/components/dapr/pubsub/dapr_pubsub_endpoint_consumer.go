// //go:build components_dapr_pubsub || components_all

package pubsub

import (
	"context"
	"strings"
	"sync/atomic"

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
		}

		c.Logger().Debugf("subscribing to %+v", sub)

		err := c.endpoint.s.AddTopicEventHandler(&sub, c.handler)
		if err != nil {
			return err
		}

		c.Logger().Debugf("subscribed to %+v", sub)

		err = c.endpoint.s.Start()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Consumer) Stop(context.Context) error {
	if c.running.CompareAndSwap(true, false) {
		if c.endpoint.s == nil {
			return nil
		}

		return c.endpoint.s.Stop()
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

	c.Logger().Infof("event - PubsubName: %s, Topic: %s, ID: %s, Content-Type: %s, Data: %s",
		e.PubsubName,
		e.Topic,
		e.ID,
		e.DataContentType,
		string(e.RawData))

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
