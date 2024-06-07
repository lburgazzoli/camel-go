package kafka

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/redpanda"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/containers"
	"github.com/testcontainers/testcontainers-go"
	"github.com/twmb/franz-go/pkg/kadm"

	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/pkg/errors"
	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	ContainerType        = "redpanda"
	DefaultDialerTimeout = 10 * time.Second
)

type Container struct {
	*redpanda.Container
}

func (c *Container) Stop(ctx context.Context) error {
	if c == nil {
		return nil
	}

	if err := c.Terminate(ctx); err != nil {
		return errors.Wrap(err, "failed to terminate container")
	}

	return nil
}

func (c *Container) Client(ctx context.Context, opts ...kgo.Opt) (*kgo.Client, error) {
	id := uuid.New()

	broker, err := c.Broker(ctx)
	if err != nil {
		return nil, err
	}

	kopts := make([]kgo.Opt, 0, len(opts)+1)
	kopts = append(kopts, opts...)
	kopts = append(kopts, kgo.SeedBrokers(broker))
	kopts = append(kopts, kgo.WithLogger(kgo.BasicLogger(os.Stdout, kgo.LogLevelInfo, func() string { return id })))
	kopts = append(kopts, kgo.Dialer(func(ctx context.Context, network string, host string) (net.Conn, error) {
		dialer := &net.Dialer{Timeout: DefaultDialerTimeout}
		return dialer.DialContext(ctx, "tcp4", host)
	}))

	return kgo.NewClient(kopts...)
}

func (c *Container) Admin(ctx context.Context) (*kadm.Client, error) {
	client, err := c.Client(ctx)
	if err != nil {
		return nil, err
	}

	return kadm.NewClient(client), nil
}

func (c *Container) Broker(ctx context.Context) (string, error) {
	return c.Container.KafkaSeedBroker(ctx)
}

func (c *Container) Properties(ctx context.Context) (map[string]any, error) {
	broker, err := c.Broker(ctx)
	if err != nil {
		return nil, err
	}

	props := map[string]any{
		"kafka.broker": broker,
	}

	return props, nil
}

func NewContainer(ctx context.Context) (*Container, error) {
	container, err := redpanda.RunContainer(
		ctx,
		testcontainers.WithLogger(containers.NewSlogLogger(ContainerType)),
		redpanda.WithAutoCreateTopics(),
	)
	if err != nil {
		return nil, err
	}

	c := Container{
		Container: container,
	}

	return &c, nil
}
