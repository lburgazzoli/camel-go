////go:build components_mqtt_v3 || components_all

package v3

import (
	"context"
	"encoding/json"

	"github.com/lburgazzoli/camel-go/pkg/cloudevents"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/pkg/errors"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
)

type Producer struct {
	processors.DefaultVerticle

	endpoint *Endpoint
	client   mqtt.Client
	tc       api.TypeConverter
}

func (p *Producer) Endpoint() api.Endpoint {
	return p.endpoint
}

func (p *Producer) Start(context.Context) error {
	p.client = p.endpoint.newClient()
	if token := p.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (p *Producer) Stop(context.Context) error {
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
func (p *Producer) publish(_ context.Context, msg api.Message) {
	payload := cloudevents.CloudEventJSON{
		SpecVersion:       "1.0",
		ID:                msg.ID(),
		Type:              msg.Type(),
		Source:            msg.Source(),
		Subject:           msg.Subject(),
		Time:              msg.Time().String(),
		DataContentType:   msg.ContentType(),
		DataContentSchema: msg.ContentSchema(),
	}

	var content []byte

	_, err := p.tc.Convert(msg.Content(), &content)
	if err != nil {
		msg.SetError(errors.Wrap(err, "error converting content to []byte"))
		return
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		msg.SetError(errors.Wrap(err, "error converting content to []byte"))
		return
	}

	t := p.client.Publish(
		p.endpoint.config.Remaining,
		p.endpoint.config.QoS,
		p.endpoint.config.Retained,
		bytes)

	t.Wait()
}
