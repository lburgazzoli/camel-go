package mqtt

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/lburgazzoli/camel-go/test/support/containers"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
)

const (
	DefaultPort          = 1883
	DefaultWebsocketPort = 9001
	DefaultVersion       = "2.0.15"
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

func NewContainer(ctx context.Context, overrideReq containers.OverrideContainerRequestOption) (*Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("docker.io/eclipse-mosquitto:%s", DefaultVersion),
		ExposedPorts: []string{fmt.Sprintf("%d", DefaultPort), fmt.Sprintf("%d", DefaultWebsocketPort)},
		Env:          map[string]string{},
		WaitingFor:   wait.ForListeningPort(nat.Port(fmt.Sprintf("%d", DefaultPort))),
		SkipReaper:   os.Getenv("TESTCONTAINERS_RYUK_DISABLED") == "true",
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
