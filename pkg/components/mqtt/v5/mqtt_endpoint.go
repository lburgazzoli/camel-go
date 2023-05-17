////go:build components_mqtt_v5 || components_all

package v5

import (
	"context"
	"net"
	"net/url"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/pkg/errors"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"

	"github.com/eclipse/paho.golang/paho"
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

func (e *Endpoint) newClient(opts ...OptionFn) (*Client, error) {

	cc := paho.ClientConfig{}

	for _, fn := range opts {
		fn(&cc)
	}

	u, err := url.Parse(e.config.Broker)
	if err != nil {
		return nil, errors.Wrapf(err, "iunvalid broker url %s", e.config.Broker)
	}

	conn, err := net.Dial(u.Scheme, u.Host)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to dial to %s", e.config.Broker)
	}

	sl := e.Logger().Sugar()

	cc.Conn = conn
	cc.OnServerDisconnect = func(disconnect *paho.Disconnect) {
		sl.Warnf("disconnected %v", disconnect)
	}
	cc.OnClientError = func(err error) {
		sl.Warnf("client error %v", err)
	}

	c := Client{
		logger: sl,
		cfg:    &e.config,
		client: paho.NewClient(cc),
	}

	return &c, nil
}
