////go:build components_mqtt_v5 || components_all

package v5

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/eclipse/paho.golang/paho"
	"github.com/pkg/errors"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

type Producer struct {
	processors.DefaultVerticle

	running  atomic.Bool
	endpoint *Endpoint
	client   *Client
	tc       api.TypeConverter
}

func (p *Producer) Endpoint() api.Endpoint {
	return p.endpoint
}

func (p *Producer) Start(ctx context.Context) error {
	if p.running.CompareAndSwap(false, true) {
		cl, err := p.endpoint.newClient()

		if err != nil {
			return err
		}

		p.client = cl

		err = p.client.Start(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Producer) Stop(ctx context.Context) error {
	if p.running.CompareAndSwap(true, false) {
		return p.client.Stop(ctx)
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

// publish produces a message that conforms to the CloudEvents binary-content-mode 1.0 spec.
func (p *Producer) publish(ctx context.Context, msg api.Message) {
	props := paho.PublishProperties{}
	props.ContentType = msg.ContentType()
	props.User.Add("ce_specversion", "1.0")

	// copy relevant attributes as ce headers
	_ = msg.EachAttribute(func(k string, v any) error {
		switch k {
		case api.MessageAttributeID:
			k = "id"
		case api.MessageAttributeTime:
			k = "time"
		case api.MessageAttributeSource:
			k = "source"
		case api.MessageAttributeContentSchema:
			k = "datacontentschema"
		default:
			return nil
		}

		err := p.setUserProperty(&props, k, v)
		if err != nil {
			msg.SetError(err)
			return nil
		}

		return nil
	})

	// copy remaining headers a standard headers
	_ = msg.EachHeader(func(k string, v any) error {
		err := p.setUserProperty(&props, k, v)
		if err != nil {
			msg.SetError(err)
			return nil
		}

		return nil
	})

	pb := &paho.Publish{
		Topic:      p.endpoint.config.Remaining,
		QoS:        p.endpoint.config.QoS,
		Payload:    []byte{},
		Properties: &props,
	}

	err := p.client.Publish(ctx, pb)
	if err != nil {
		msg.SetError(errors.Wrapf(
			err,
			"error while publishing to topic '%s' on server '%s'",
			p.endpoint.config.Remaining,
			p.endpoint.config.Broker),
		)
	}
}

func (p *Producer) setUserProperty(properties *paho.PublishProperties, k string, v any) error {
	var val string

	ok, err := p.tc.Convert(v, &val)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("unable to convert value for header %s", k)
	}

	properties.User.Add(k, val)

	return nil
}
