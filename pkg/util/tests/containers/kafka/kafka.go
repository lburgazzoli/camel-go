package kafka

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/redpanda"

	"github.com/docker/go-connections/nat"
	"github.com/lburgazzoli/camel-go/pkg/util/tests/containers"
	"github.com/testcontainers/testcontainers-go"
	"github.com/twmb/franz-go/pkg/kadm"

	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/pkg/errors"
	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	DefaultImageName            = "docker.io/redpandadata/redpanda"
	DefaultImageVersion         = "v24.1.7"
	ContainerType               = "redpanda"
	ContainerEntrypointFile     = "/entrypoint-tc.sh"
	RedPandaDir                 = "/etc/redpanda"
	RedPandaBootstrapConfigFile = ".bootstrap.yaml"
	RedPandaBConfigFile         = "redpanda.yaml"

	DefaultPort            = 9092
	DefaultAdminPort       = 9644
	DefaultLogPollInterval = 100 * time.Millisecond
	DefaultDialerTimeout   = 10 * time.Second
)

type Container struct {
	*redpanda.Container
}

func (c *Container) Stop(ctx context.Context) error {
	if c == nil {
		return nil
	}

	if err := c.StopLogProducer(); err != nil {
		return errors.Wrap(err, "failed to  stop log producers")
	}

	if err := c.Terminate(ctx); err != nil {
		return errors.Wrap(err, "failed to terminate container")
	}

	return nil
}

func (c *Container) Client(ctx context.Context, opts ...kgo.Opt) (*kgo.Client, error) {
	id := uuid.New()

	host, err := c.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := c.MappedPort(ctx, nat.Port(fmt.Sprintf("%d/tcp", DefaultPort)))
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped Kafka port: %w", err)
	}

	kopts := make([]kgo.Opt, 0, len(opts)+1)
	kopts = append(kopts, opts...)
	kopts = append(kopts, kgo.SeedBrokers(host+":"+port.Port()))
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
	host, err := c.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := c.MappedPort(ctx, nat.Port(fmt.Sprintf("%d/tcp", DefaultPort)))
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped Kafka port: %w", err)
	}

	broker, err := c.Broker(ctx)
	if err != nil {
		return nil, err
	}

	props := map[string]any{
		"kafka.broker": broker,
		"kafka.host":   host,
		"kafka.port":   port.Int(),
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

	if err := container.StartLogProducer(ctx); err != nil {
		return nil, err
	}

	c := Container{
		Container: container,
	}

	return &c, nil
}
