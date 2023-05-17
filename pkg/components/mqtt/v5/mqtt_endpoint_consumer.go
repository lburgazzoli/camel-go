////go:build components_mqtt_v5 || components_all

package v5

import (
	"context"
	"strconv"
	"sync/atomic"

	"github.com/pkg/errors"

	"github.com/eclipse/paho.golang/paho"
	"github.com/lburgazzoli/camel-go/pkg/components"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
)

type Consumer struct {
	components.DefaultConsumer

	running  atomic.Bool
	endpoint *Endpoint
	client   *Client
}

func (c *Consumer) Endpoint() camel.Endpoint {
	return c.endpoint
}

func (c *Consumer) Start(ctx context.Context) error {
	if c.running.CompareAndSwap(false, true) {
		cl, err := c.endpoint.newClient(
			WithSingleHandlerRouter(c.handler),
		)

		if err != nil {
			return err
		}

		c.client = cl

		err = c.client.Start(ctx)
		if err != nil {
			return err
		}

		err = c.client.Subscribe(ctx, c.endpoint.config.Remaining)
		if err != nil {
			return errors.Wrapf(err, "failure subscribing to topic %s", c.endpoint.config.Remaining)
		}
	}

	return nil
}

func (c *Consumer) Stop(ctx context.Context) error {
	if c.running.CompareAndSwap(true, false) {
		_ = c.client.Stop(ctx)
		c.client = nil
	}

	return nil
}

func (c *Consumer) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		err := c.Start(context.Background())
		if err != nil {
			panic(err)
		}
	case *actor.Stopping:
		err := c.Stop(context.Background())
		if err != nil {
			panic(err)
		}
	case camel.Message:
		// ignore message,
		// TODO: may be used for ack
		break
	}
}

func (c *Consumer) handler(pub *paho.Publish) {
	c.Logger().Infof("handling publish %v", pub)

	component := c.endpoint.Component()
	camelCtx := component.Context()

	m := camelCtx.NewMessage()

	m.SetSubject(pub.Topic)
	m.SetType("event")
	m.SetSource(component.Scheme())
	m.SetContent(pub.Payload)

	if pub.Properties != nil {
		m.SetContentType(pub.Properties.ContentType)
	}

	_ = m.SetAttribute(AttributeMqttMessageID, strconv.FormatUint(uint64(pub.PacketID), 10))
	_ = m.SetAttribute(AttributeMqttMessageRetained, pub.Retain)
	_ = m.SetAttribute(AttributeMqttMessageQUOS, pub.QoS)

	if err := component.Context().SendTo(c.Target(), m); err != nil {
		panic(err)
	}
}
