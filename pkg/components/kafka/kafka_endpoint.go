////go:build components_kafka || components_all

package kafka

import (
	"context"
	"crypto/tls"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"github.com/twmb/franz-go/plugin/kzap"
	"net"
	"strings"
	"time"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
)

type Endpoint struct {
	config Config
	components.DefaultEndpoint
}

func (e *Endpoint) Start(context.Context) error {
	return nil
}

func (e *Endpoint) Stop(context.Context) error {
	return nil
}

func (e *Endpoint) Producer() (api.Producer, error) {
	c := Producer{
		DefaultVerticle: processors.NewDefaultVerticle(),
		endpoint:        e,
		tc:              e.Context().TypeConverter(),
	}

	return &c, nil
}

func (e *Endpoint) Consumer(pid *actor.PID) (api.Consumer, error) {
	c := Consumer{
		DefaultConsumer: components.NewDefaultConsumer(e, pid),
		endpoint:        e,
	}

	return &c, nil
}

func (e *Endpoint) newClient() (*kgo.Client, error) {
	opts := make([]kgo.Opt, 0)
	opts = append(opts, kgo.SeedBrokers(strings.Split(e.config.Brokers, ",")...))

	if e.config.User != "" && e.config.Password != "" {
		tlsDialer := &tls.Dialer{NetDialer: &net.Dialer{Timeout: 10 * time.Second}}
		authMechanism := plain.Auth{User: e.config.User, Pass: e.config.Password}.AsMechanism()

		opts = append(opts, kgo.SASL(authMechanism))
		opts = append(opts, kgo.Dialer(tlsDialer.DialContext))
		opts = append(opts, kgo.WithLogger(kzap.New(e.Logger())))
	}

	return kgo.NewClient(opts...)
}
