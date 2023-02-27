package kafka

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/twmb/franz-go/pkg/kadm"

	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/lburgazzoli/camel-go/test/support/containers"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	DefaultPort    = 9092
	DefaultVersion = "v22.3.13"
)

type Request struct {
	testcontainers.ContainerRequest
}

type Container struct {
	testcontainers.Container
}

func (c *Container) Stop(ctx context.Context) error {
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
		return nil, err
	}

	kopts := make([]kgo.Opt, 0, len(opts)+1)
	kopts = append(kopts, opts...)
	kopts = append(kopts, kgo.SeedBrokers(host+":"+strconv.Itoa(DefaultPort)))
	kopts = append(kopts, kgo.WithLogger(kgo.BasicLogger(os.Stdout, kgo.LogLevelInfo, func() string { return id })))

	return kgo.NewClient(kopts...)
}

func (c *Container) Admin(ctx context.Context) (*kadm.Client, error) {

	client, err := c.Client(ctx)
	if err != nil {
		return nil, err
	}

	return kadm.NewClient(client), nil
}

func NewContainer(ctx context.Context, overrideReq containers.OverrideContainerRequestOption) (*Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("docker.io/redpandadata/redpanda:%s", DefaultVersion),
		ExposedPorts: []string{fmt.Sprintf("%d:%d", DefaultPort, DefaultPort)},
		Env:          map[string]string{},
		WaitingFor:   wait.ForLog("Started Kafka API server"),
		Cmd:          []string{"redpanda", "start", "--mode dev-container"},
	}
	// }

	kafkaRequest := Request{
		ContainerRequest: req,
	}

	if overrideReq != nil {
		merged := overrideReq(kafkaRequest.ContainerRequest)
		kafkaRequest.ContainerRequest = merged
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: kafkaRequest.ContainerRequest,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	c := Container{
		Container: container,
	}

	c.FollowOutput(&containers.SysOutLogConsumer{})

	if err := container.StartLogProducer(ctx); err != nil {
		return nil, err
	}

	return &c, nil
}
