////go:build components_all || components_timer || steps_all || steps_process

package engine

import (
	"context"
	"fmt"
	"strings"

	"github.com/lburgazzoli/camel-go/pkg/components/timer"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	cameltest "github.com/lburgazzoli/camel-go/test/support"

	"testing"
	"time"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/endpoint"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/from"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/process"
	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	wg := make(chan api.Message)

	c := core.NewContext()
	assert.NotNil(t, c)

	c.Registry().Set("consumer", func(message api.Message) {
		wg <- message
	})

	f := from.From{
		Endpoint: endpoint.Endpoint{
			URI: "timer:foo",
			Parameters: map[string]interface{}{
				"interval": 1 * time.Second,
			},
		},
	}

	f.Next(cameltest.Reify(t, c, &process.Process{Ref: "consumer"}))

	fromPid, err := f.Reify(c)
	assert.Nil(t, err)
	assert.NotNil(t, fromPid)

	select {
	case msg := <-wg:
		a, ok := msg.Annotation(timer.AnnotationTimerFiredCount)
		assert.True(t, ok)
		assert.Equal(t, uint64(1), a)
	case <-time.After(5 * time.Second):
		assert.Fail(t, "timeout")
	}
}

const simpleRoute = `
- route:
    from:
      uri: "timer:foo"
      steps:
        - process:
            ref: "consumer-1"
        - process:
            ref: "consumer-2"
`

func TestSimpleYAML(t *testing.T) {
	content := uuid.New()
	wg := make(chan api.Message)

	c := core.NewContext()
	assert.NotNil(t, c)

	c.Registry().Set("consumer-1", func(message api.Message) {
		message.SetContent(content)
	})
	c.Registry().Set("consumer-2", func(message api.Message) {
		wg <- message
	})

	err := c.LoadRoutes(strings.NewReader(simpleRoute))
	assert.Nil(t, err)

	select {
	case msg := <-wg:
		a, ok := msg.Annotation(timer.AnnotationTimerFiredCount)
		assert.True(t, ok)
		assert.Equal(t, uint64(1), a)
		assert.Equal(t, content, msg.Content())
	case <-time.After(5 * time.Second):
		assert.Fail(t, "timeout")
	}
}

const _ = `
- route:
    from:
      uri: "timer:foo"
      steps:
        - process:
            ref: "consumer-1"
        - process:
            ref: "consumer-2"
`

func TestSimpleKafka(t *testing.T) {
	t.Skip("TODO")

	ctx := context.Background()

	_ = strings.Join(
		[]string{
			"OUTSIDE://0.0.0.0:9092",
			"PLAINTEXT://0.0.0.0:9092",
		},
		",")

	_ = strings.Join(
		[]string{
			"OUTSIDE://localhost:9092",
			"PLAINTEXT://localhost:9092",
		},
		",")

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "docker.io/redpandadata/redpanda:v22.3.13",
			ExposedPorts: []string{"9092:9092"},
			WaitingFor:   wait.ForLog("Started Kafka API server"),
			Cmd:          []string{"redpanda", "start", "--mode dev-container"},
		},
		Started: true,
	})

	if err != nil {
		t.Error(err)
	}

	defer func() {
		if err := container.StopLogProducer(); err != nil {
			t.Fatalf("failed to stop log producers: %s", err.Error())
		}
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()

	container.FollowOutput(&TestLogConsumer{})
	assert.Nil(t, container.StartLogProducer(ctx))
	assert.Nil(t, container.Start(ctx))

	host, err := container.Host(ctx)
	assert.Nil(t, err)

	port, err := container.MappedPort(ctx, "9092")
	assert.Nil(t, err)

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(host + ":" + port.Port()),
	)

	assert.Nil(t, err)

	defer cl.Close()

	ac := kadm.NewClient(cl)

	tp, err := ac.CreateTopic(ctx, 3, 1, nil, "foo")
	assert.Nil(t, err)
	assert.Nil(t, tp.Err)

	assert.Nil(t, cl.ProduceSync(
		ctx,
		&kgo.Record{
			Topic: "foo",
			Value: []byte("bar"),
		}).FirstErr())

	for {
		fetches := cl.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			// All errors are retried internally when fetching, but non-retriable errors are
			// returned from polls so that users can notice and take action.
			panic(fmt.Sprint(errs))
		}

		// We can iterate through a record iterator...
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			fmt.Println(string(record.Value), "from an iterator!")
		}
	}
}

type TestLogConsumer struct {
}

func (g *TestLogConsumer) Accept(l testcontainers.Log) {
	fmt.Print(string(l.Content))
}
