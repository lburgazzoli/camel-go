////go:build components_all || components_timer || steps_all || steps_process

package engine

import (
	"context"

	"github.com/lburgazzoli/camel-go/test/support/containers"
	"github.com/lburgazzoli/camel-go/test/support/containers/kafka"
	. "github.com/onsi/gomega"

	"strings"

	_ "github.com/lburgazzoli/camel-go/pkg/components/kafka"
	"github.com/lburgazzoli/camel-go/pkg/components/timer"
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
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/to"
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

const simpleKafka = `
- route:
    from:
      uri: "timer:foo"
      steps:
        - process:
            ref: "consumer-1"
        - to:
            uri: "kafka:foo"
            parameters:
              brokers: "localhost:9092"
              topics: "foo"
`

func TestSimpleKafka(t *testing.T) {
	content := uuid.New()
	ctx := context.Background()

	container, err := kafka.NewContainer(ctx, containers.NoopOverrideContainerRequest)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		if err := container.Stop(ctx); err != nil {
			t.Fatal(err.Error())
		}
	}()

	assert.Nil(t, container.Start(ctx))

	cl, err := container.Client(
		ctx,
		kgo.ConsumeTopics("foo"),
		kgo.ConsumerGroup(uuid.New()),
	)

	assert.Nil(t, err)

	defer cl.Close()

	ac := kadm.NewClient(cl)

	tp, err := ac.CreateTopic(ctx, 3, 1, nil, "foo")
	assert.Nil(t, err)
	assert.Nil(t, tp.Err)

	c := core.NewContext()
	assert.NotNil(t, c)

	c.Registry().Set("consumer-1", func(message api.Message) {
		message.SetContent(content)
	})

	err = c.LoadRoutes(strings.NewReader(simpleKafka))
	assert.Nil(t, err)

	RegisterTestingT(t)

	Eventually(func(g Gomega) {
		f := cl.PollFetches(ctx)

		Expect(f.Errors()).To(BeEmpty())
		Expect(f.NumRecords()).To(Equal(1))
		Expect(string(f.Records()[0].Value)).To(Equal(content))
	}).Should(Succeed())
}
