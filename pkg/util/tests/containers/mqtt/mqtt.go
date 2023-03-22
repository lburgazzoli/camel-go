package mqtt

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/containers"

	"github.com/docker/go-connections/nat"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	paho "github.com/eclipse/paho.mqtt.golang"
)

const (
	DefaultPort          = 1883
	DefaultWebsocketPort = 9001
	DefaultVersion       = "2.0.15"
)

type RequestFn func(*Request) *Request

type Request struct {
	testcontainers.ContainerRequest
	Config string
}

func WithConfig(path string) RequestFn {
	return func(request *Request) *Request {
		request.Config = path
		return request
	}
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

func (c *Container) Client(ctx context.Context) (paho.Client, error) {
	broker, err := c.Broker(ctx)
	if err != nil {
		return nil, err
	}

	opts := paho.NewClientOptions()
	opts = opts.AddBroker(broker)
	opts = opts.SetClientID(uuid.New())
	opts = opts.SetKeepAlive(2 * time.Second)
	opts = opts.SetPingTimeout(1 * time.Second)

	// opts.ConnectRetry = true
	// opts.AutoReconnect = true

	// Log events
	opts.OnConnectionLost = func(cl paho.Client, err error) {
		fmt.Println("connection lost")
	}
	opts.OnConnect = func(paho.Client) {
		fmt.Println("connection established")
	}
	opts.OnReconnecting = func(paho.Client, *paho.ClientOptions) {
		fmt.Println("attempting to reconnect")
	}

	client := paho.NewClient(opts)

	// TODO: must not blocking probably
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}

//nolint:misspell
func NewContainer(ctx context.Context, opts ...RequestFn) (*Container, error) {
	req := &Request{
		ContainerRequest: testcontainers.ContainerRequest{
			SkipReaper: os.Getenv("TESTCONTAINERS_RYUK_DISABLED") == "true",
			Image:      fmt.Sprintf("docker.io/eclipse-mosquitto:%s", DefaultVersion),
			Env:        map[string]string{},
			ExposedPorts: []string{
				fmt.Sprintf("%d", DefaultPort),
			},
			WaitingFor: wait.ForAll(
				wait.ForExposedPort(),
			),
		},
	}

	for i := range opts {
		req = opts[i](req)
	}

	if req.Config != "" {
		req.ContainerRequest.Files = append(req.ContainerRequest.Files, testcontainers.ContainerFile{
			HostFilePath:      req.Config,
			ContainerFilePath: "/mosquitto/config/mosquitto.conf",
			FileMode:          664,
		})
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req.ContainerRequest,
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
