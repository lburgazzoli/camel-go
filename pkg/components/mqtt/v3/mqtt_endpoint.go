////go:build components_mqtt_v3 || components_all

package v3

import (
	"context"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
)

type Endpoint struct {
	components.DefaultEndpoint

	config Config
}

func (e *Endpoint) Start(context.Context) error {
	return nil
}

func (e *Endpoint) Stop(context.Context) error {
	return nil
}

func (e *Endpoint) Consumer(pid *actor.PID) (api.Consumer, error) {
	c := Consumer{
		DefaultConsumer: components.NewDefaultConsumer(e, pid),
		endpoint:        e,
	}

	return &c, nil
}

func (e *Endpoint) Producer() (api.Producer, error) {
	p := Producer{
		DefaultVerticle: processors.NewDefaultVerticle(),
		endpoint:        e,
		tc:              e.Context().TypeConverter(),
	}

	return &p, nil
}

func (e *Endpoint) newClient(optFns ...OptionFn) mqtt.Client {
	cid := e.config.ClientID
	if cid == "" {
		cid = uuid.New()
	}

	opts := mqtt.NewClientOptions()
	opts = opts.SetClientID(cid)
	opts = opts.SetKeepAlive(2 * time.Second)
	opts = opts.SetPingTimeout(1 * time.Second)

	for _, fn := range optFns {
		fn(opts)
	}

	for _, broker := range strings.Split(e.config.Broker, ",") {
		if broker == "" {
			continue
		}

		opts = opts.AddBroker(broker)
	}

	sl := e.Logger().Sugar()

	// Log events
	opts.OnConnectionLost = func(cl mqtt.Client, err error) {
		sl.Warnf("connection lost (error: %s)", err.Error())
	}
	opts.OnConnect = func(cl mqtt.Client) {
		sl.Info("connection established")
	}
	opts.OnReconnecting = func(mqtt.Client, *mqtt.ClientOptions) {
		sl.Info("attempting to reconnect")
	}

	return mqtt.NewClient(opts)
}
