// //go:build components_kafka || components_all

package kafka

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net"
	"strings"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"

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
		DefaultProducer: components.NewDefaultProducer(e),
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

func (e *Endpoint) newClient(additionalOpts ...kgo.Opt) (*kgo.Client, error) {
	opts := make([]kgo.Opt, 0)
	opts = append(opts, kgo.SeedBrokers(strings.Split(e.config.Brokers, ",")...))
	opts = append(opts, kgo.WithLogger(&klog{delegate: e.Logger().With(slog.String("subsystem", "kafka"))}))

	dialer := &net.Dialer{Timeout: 10 * time.Second}

	if e.config.Username != "" && e.config.Password != "" {
		tlsDialer := &tls.Dialer{NetDialer: dialer}
		authMechanism := plain.Auth{User: e.config.Username, Pass: e.config.Password}.AsMechanism()

		opts = append(opts, kgo.SASL(authMechanism))
		opts = append(opts, kgo.Dialer(func(ctx context.Context, network string, host string) (net.Conn, error) {
			n := network
			if e.config.ForceIPV4 {
				n = "tcp4"
			}

			return tlsDialer.DialContext(ctx, n, host)
		}))
	} else {
		opts = append(opts, kgo.Dialer(func(ctx context.Context, network string, host string) (net.Conn, error) {
			n := network
			if e.config.ForceIPV4 {
				n = "tcp4"
			}

			return dialer.DialContext(ctx, n, host)
		}))
	}

	opts = append(opts, additionalOpts...)

	return kgo.NewClient(opts...)
}
