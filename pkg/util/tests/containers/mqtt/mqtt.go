package mqtt

import (
	"context"
	"log/slog"
	"net"
	"strconv"
	"time"

	"github.com/rs/xid"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/containers"

	"github.com/docker/go-connections/nat"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	paho "github.com/eclipse/paho.mqtt.golang"
)

const (
	DefaultPort        = 1883
	DefaultVersion     = "2.0.15"
	DefaultKeepAlive   = 2 * time.Second
	DefaultPingTimeout = 1 * time.Second
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
	if c == nil {
		return nil
	}

	if err := c.Terminate(ctx); err != nil {
		return errors.Wrap(err, "failed to terminate container")
	}

	return nil
}

func (c *Container) Broker(ctx context.Context) (string, error) {
	address, err := c.Address(ctx)
	if err != nil {
		return "", err
	}

	return "tcp://" + address, nil
}

func (c *Container) Address(ctx context.Context) (string, error) {
	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}

	port, err := c.MappedPort(ctx, nat.Port(strconv.Itoa(DefaultPort)))
	if err != nil {
		return "", err
	}

	return net.JoinHostPort(host, port.Port()), nil
}

func (c *Container) Client(ctx context.Context) (paho.Client, error) {
	broker, err := c.Broker(ctx)
	if err != nil {
		return nil, err
	}

	name, err := c.Name(ctx)
	if err != nil {
		return nil, err
	}

	opts := paho.NewClientOptions()
	opts = opts.AddBroker(broker)
	opts = opts.SetClientID(uuid.New())
	opts = opts.SetKeepAlive(DefaultKeepAlive)
	opts = opts.SetPingTimeout(DefaultPingTimeout)

	// opts.ConnectRetry = true
	// opts.AutoReconnect = true

	// Log events
	opts.OnConnectionLost = func(cl paho.Client, err error) {
		containers.Log.Info("connection lost", slog.String("container", name), slog.String("error", err.Error()))
	}
	opts.OnConnect = func(paho.Client) {
		containers.Log.Info("connection established", slog.String("container", name))
	}
	opts.OnReconnecting = func(paho.Client, *paho.ClientOptions) {
		containers.Log.Info("attempting to reconnect", slog.String("container", name))
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
			Name:  "eclipse-mosquitto" + xid.New().String(),
			Image: "docker.io/eclipse-mosquitto:" + DefaultVersion,
			Env:   map[string]string{},
			ExposedPorts: []string{
				strconv.Itoa(DefaultPort),
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
			FileMode:          int64(containers.FileModeShared),
		})
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req.ContainerRequest,
		Started:          true,
		Logger:           containers.NewSlogLogger("eclipse-mosquitto"),
	})
	if err != nil {
		return nil, err
	}

	c := Container{
		Container: container,
	}

	return &c, nil
}
