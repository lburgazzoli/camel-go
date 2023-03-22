////go:build components_mqtt || components_all

package mqtt

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/asynkron/protoactor-go/actor"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

type Consumer struct {
	api.WithOutputs

	id       string
	endpoint *Endpoint
	client   mqtt.Client
	logger   *zap.SugaredLogger
}

func (c *Consumer) Endpoint() api.Endpoint {
	return c.endpoint
}

func (c *Consumer) ID() string {
	return c.id
}

func (c *Consumer) Start(context.Context) error {
	cid := c.endpoint.config.ClientID
	if cid == "" {
		cid = uuid.New()
	}

	opts := mqtt.NewClientOptions()
	opts = opts.SetClientID(cid)
	opts = opts.SetKeepAlive(2 * time.Second)
	opts = opts.SetDefaultPublishHandler(c.handler)
	opts = opts.SetPingTimeout(1 * time.Second)

	for _, broker := range strings.Split(c.endpoint.config.Brokers, ",") {
		if broker == "" {
			continue
		}

		opts = opts.AddBroker(broker)
	}

	// Log events
	opts.OnConnectionLost = func(cl mqtt.Client, err error) {
		c.logger.Warnf("connection lost (error: %s)", err.Error())
	}
	opts.OnConnect = func(cl mqtt.Client) {
		c.logger.Info("connection established")
	}
	opts.OnReconnecting = func(mqtt.Client, *mqtt.ClientOptions) {
		c.logger.Info("attempting to reconnect")
	}

	c.client = mqtt.NewClient(opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	c.logger.Infof("subscribing to %s", c.endpoint.config.Remaining)

	token := c.client.Subscribe(c.endpoint.config.Remaining, 0, c.handler)
	token.Wait()

	if token.Error() != nil {
		return token.Error()
	}

	c.logger.Infof("subscribed to %s", c.endpoint.config.Remaining)

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
	}
}

func (c *Consumer) handler(_ mqtt.Client, msg mqtt.Message) {
	c.logger.Infof("handing message %v", msg)

	m, err := message.New()
	if err != nil {
		panic(err)
	}

	component := c.endpoint.Component()

	_ = m.SetSubject(msg.Topic())
	_ = m.SetType("mqtt.publish")
	_ = m.SetSource(component.Scheme())

	m.SetAnnotation(AnnotationMqttMessageID, strconv.FormatUint(uint64(msg.MessageID()), 10))
	m.SetContent(msg.Payload())

	for _, o := range c.Outputs() {
		if err := component.Context().SendTo(o, m); err != nil {
			panic(err)
		}
	}

	msg.Ack()
}
