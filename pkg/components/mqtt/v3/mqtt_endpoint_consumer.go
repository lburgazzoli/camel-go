// //go:build components_mqtt_v3 || components_all

package v3

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/lburgazzoli/camel-go/pkg/components"

	"github.com/asynkron/protoactor-go/actor"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
)

type Consumer struct {
	components.DefaultConsumer

	endpoint *Endpoint
	client   mqtt.Client
}

func (c *Consumer) Endpoint() camel.Endpoint {
	return c.endpoint
}

func (c *Consumer) Start(context.Context) error {
	c.client = c.endpoint.newClient(func(opts *mqtt.ClientOptions) {
		opts.SetDefaultPublishHandler(c.handler)
	})

	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	c.Logger().Info("subscribing", slog.String("topic", c.endpoint.config.Remaining))

	token := c.client.Subscribe(c.endpoint.config.Remaining, 0, c.handler)
	token.Wait()

	if token.Error() != nil {
		return token.Error()
	}

	c.Logger().Info("subscribed", slog.String("topic", c.endpoint.config.Remaining))

	return nil
}

func (c *Consumer) Stop(context.Context) error {
	if token := c.client.Unsubscribe(c.endpoint.config.Remaining); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	if c.client != nil {
		c.client.Disconnect(250)
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
		// TODO: may be used for transactions
		break
	}
}

func (c *Consumer) handler(_ mqtt.Client, msg mqtt.Message) {
	c.Logger().Debug("message received", slog.Group(
		"message",
		slog.Uint64("id", uint64(msg.MessageID())),
		slog.String("topic", msg.Topic()),
		slog.Bool("retained", msg.Retained()),
	))

	component := c.endpoint.Component()
	camelCtx := component.Context()

	m := camelCtx.NewMessage()

	m.SetSubject(msg.Topic())
	m.SetType("event")
	m.SetSource(component.Scheme())
	m.SetContent(msg.Payload())

	m.SetAttribute(AttributeMqttMessageID, strconv.FormatUint(uint64(msg.MessageID()), 10))
	m.SetAttribute(AttributeMqttMessageRetained, msg.Retained())
	m.SetAttribute(AttributeMqttMessageDuplicate, msg.Duplicate())
	m.SetAttribute(AttributeMqttMessageQUOS, msg.Qos())

	if err := component.Context().SendTo(c.Target(), m); err != nil {
		panic(err)
	}

	msg.Ack()
}
