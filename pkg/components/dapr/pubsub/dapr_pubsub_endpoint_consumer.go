////go:build components_dapr_pubsub || components_all

package pubsub

import (
	"context"
	"sync/atomic"

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
		sub := common.Subscription{
			PubsubName: c.pubsubName,
			Topic:      c.topicName,
			Route:      "/" + c.endpoint.ID() + "/" + c.ID(),
		}

		err := c.endpoint.s.AddTopicEventHandler(&sub, c.handler)
		if err != nil {
			return err
		}

		return c.endpoint.s.Start()
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

	m := camelCtx.NewMessage()

	m.SetType(e.Type)
	m.SetSource(e.Source)
	m.SetSubject(e.Subject)
	m.SetContentType(e.DataContentType)

	_ = m.SetAttribute(AttributeEventID, e.ID)
	_ = m.SetAttribute(AttributePubSubName, e.PubsubName)
	_ = m.SetAttribute(AttributePubSubTopic, e.Topic)

	switch {
	case e.Data != nil:
		m.SetContent(e.Data)
	case e.RawData != nil:
		m.SetContent(e.RawData)
	}

	if err := camelCtx.SendTo(c.Target(), m); err != nil {
		panic(err)
	}

	return false, nil
}
