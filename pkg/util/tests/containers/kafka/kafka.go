package kafka

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"text/template"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/lburgazzoli/camel-go/pkg/util/tests/containers"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/twmb/franz-go/pkg/kadm"

	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	DefaultPort                 = 9092
	DefaultAdminPort            = 9644
	DefaultVersion              = "v23.2.12"
	ContainerEntrypointFile     = "/entrypoint-tc.sh"
	RedPandaDir                 = "/etc/redpanda"
	RedPandaBootstrapConfigFile = ".bootstrap.yaml"
	RedPandaBConfigFile         = "redpanda.yaml"
)

const contentEntrypoint = `#!/usr/bin/env bash

# Wait for testcontainer's injected redpanda config with the port only known after docker start
until grep -q "# Injected by testcontainers" "/etc/redpanda/redpanda.yaml"
do
  sleep 0.1
done
exec /entrypoint.sh $@
`

const contentBootstrap = `
auto_create_topics_enabled: true
`
const contentRedpanda = `
# Injected by testcontainers
redpanda:
  admin:
    address: 0.0.0.0
    port: 9644

  kafka_api:
    - address: 0.0.0.0
      name: external
      port: 9092
      authentication_method: none

    - address: 0.0.0.0
      name: internal
      port: 9093
      authentication_method: none

  advertised_kafka_api:
    - address: {{ .KafkaAPI.AdvertisedHost }}
      name: external
      port: {{ .KafkaAPI.AdvertisedPort }}
    - address: 127.0.0.1
      name: internal
      port: 9093

auto_create_topics_enabled: true
`

type RequestFn func(*Request) *Request

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
	return c.PortEndpoint(ctx, nat.Port(fmt.Sprintf("%d/tcp", DefaultPort)), "")
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
func NewContainer(ctx context.Context, opts ...RequestFn) (*Container, error) {
	tmpDir, err := os.MkdirTemp("", "redpanda")
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	pathEntrypoint := path.Join(tmpDir, ContainerEntrypointFile)
	pathBootstrap := path.Join(tmpDir, RedPandaBootstrapConfigFile)

	if err := os.WriteFile(pathEntrypoint, []byte(contentEntrypoint), 0o600); err != nil {
		return nil, fmt.Errorf("failed to create entrypoint file: %w", err)
	}
	if err := os.WriteFile(pathBootstrap, []byte(contentBootstrap), 0o600); err != nil {
		return nil, fmt.Errorf("failed to create entrypoint file: %w", err)
	}

	req := &Request{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:  "redpandadata-" + uuid.New(),
			Image: fmt.Sprintf("docker.io/redpandadata/redpanda:%s", DefaultVersion),
			User:  "root:root",
			Env:   map[string]string{},
			Entrypoint: []string{
				ContainerEntrypointFile,
			},
			Cmd: []string{
				"redpanda",
				"start",
				"--mode=dev-container",
				"--smp=1",
				"--memory=1G"},
			ExposedPorts: []string{
				strconv.Itoa(DefaultPort),
				strconv.Itoa(DefaultAdminPort),
			},
			Files: []testcontainers.ContainerFile{{
				HostFilePath:      pathEntrypoint,
				ContainerFilePath: ContainerEntrypointFile,
				FileMode:          700,
			}, {
				HostFilePath:      pathBootstrap,
				ContainerFilePath: filepath.Join(RedPandaDir, RedPandaBootstrapConfigFile),
				FileMode:          600,
			}},
		},
	}

	for i := range opts {
		req = opts[i](req)
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req.ContainerRequest,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	container.FollowOutput(containers.NewSlogLogConsumer(&req.ContainerRequest))

	if err := container.StartLogProducer(ctx); err != nil {
		return nil, err
	}

	contentConfig, err := generateRedPandaConfiguration(ctx, container)
	if err != nil {
		return nil, err
	}

	err = container.CopyToContainer(ctx, contentConfig, filepath.Join(RedPandaDir, RedPandaBConfigFile), 600)
	if err != nil {
		return nil, fmt.Errorf("failed to copy redpanda.yaml into container: %w", err)
	}

	c := Container{
		Container: container,
	}

	err = wait.ForLog("Successfully started Redpanda!").
		WithPollInterval(100*time.Millisecond).
		WaitUntilReady(ctx, c.Container)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func generateRedPandaConfiguration(ctx context.Context, container testcontainers.Container) ([]byte, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, nat.Port(fmt.Sprintf("%d/tcp", DefaultPort)))
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped Kafka port: %w", err)
	}

	params := map[string]any{
		"KafkaAPI": map[string]any{
			"AdvertisedHost": host,
			"AdvertisedPort": port.Int(),
		},
	}

	tpl, err := template.New(RedPandaBConfigFile).Parse(contentRedpanda)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redpanda config file template: %w", err)
	}

	var buffer bytes.Buffer
	if err := tpl.Execute(&buffer, params); err != nil {
		return nil, fmt.Errorf("failed to render redpanda node config template: %w", err)
	}

	return buffer.Bytes(), nil
}
