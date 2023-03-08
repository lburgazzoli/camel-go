package mqtt

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/go-connections/nat"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/lburgazzoli/camel-go/test/support/containers"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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

func (c *Container) Broker(ctx context.Context) (string, error) {
	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}

	port, err := c.MappedPort(ctx, nat.Port(fmt.Sprintf("%d", DefaultPort)))
	if err != nil {
		return "", err
	}

	return "tcp://" + net.JoinHostPort(host, port.Port()), nil
}

func (c *Container) Client(ctx context.Context) (mqtt.Client, error) {
	//_, err := c.Broker(ctx)
	//if err != nil {
	//	return nil, err
	//}

	opts := mqtt.NewClientOptions()
	opts = opts.AddBroker("tcp://test.mosquitto.org:1883")
	opts = opts.SetClientID(uuid.New())
	opts = opts.SetKeepAlive(2 * time.Second)
	opts = opts.SetPingTimeout(1 * time.Second)

	// opts.ConnectRetry = true
	// opts.AutoReconnect = true

	// Log events
	opts.OnConnectionLost = func(cl mqtt.Client, err error) {
		fmt.Println("connection lost")
	}
	opts.OnConnect = func(mqtt.Client) {
		fmt.Println("connection established")
	}
	opts.OnReconnecting = func(mqtt.Client, *mqtt.ClientOptions) {
		fmt.Println("attempting to reconnect")
	}

	client := mqtt.NewClient(opts)

	// TODO: must not block probably
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}

//nolint:misspell
func NewContainer(ctx context.Context, overrideReq containers.OverrideContainerRequestOption) (*Container, error) {
	conf, err := filepath.Abs("../../etc/support/mqtt/mosquitto.conf")
	if err != nil {
		return nil, err
	}

	req := testcontainers.ContainerRequest{
		SkipReaper: os.Getenv("TESTCONTAINERS_RYUK_DISABLED") == "true",
		Image:      fmt.Sprintf("docker.io/eclipse-mosquitto:%s", DefaultVersion),
		Env:        map[string]string{},
		ExposedPorts: []string{
			fmt.Sprintf("%d", DefaultPort),
		},
		WaitingFor: wait.ForAll(
			wait.ForExposedPort(),
		),
		Files: []testcontainers.ContainerFile{{
			HostFilePath:      conf,
			ContainerFilePath: "/mosquitto/config/mosquitto.conf",
			FileMode:          664,
		}},
	}

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
