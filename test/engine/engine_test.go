////go:build components_all || components_timer || steps_all || steps_process

package engine

import (
	"context"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
	"github.com/lburgazzoli/camel-go/test/support/containers"
	"github.com/lburgazzoli/camel-go/test/support/containers/kafka"
	. "github.com/onsi/gomega"

	"strings"

	_ "github.com/lburgazzoli/camel-go/pkg/components/kafka"
	"github.com/lburgazzoli/camel-go/pkg/components/timer"
	_ "github.com/lburgazzoli/camel-go/pkg/components/wasm"
	"github.com/twmb/franz-go/pkg/kgo"

	"testing"
	"time"

	"github.com/lburgazzoli/camel-go/pkg/util/uuid"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/endpoint"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/from"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/process"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/to"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/transform"
	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {

	support.Run(t, "run", func(t *testing.T, ctx context.Context, c camel.Context) {
		t.Helper()

		wg := make(chan camel.Message)

		c.Registry().Set("consumer", func(message camel.Message) {
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

		p := process.Process{Identity: uuid.New(), Ref: "consumer"}
		id, err := p.Reify(ctx, c)
		assert.Nil(t, err)
		assert.NotNil(t, id)

		f.Next(id)

		fromPid, err := f.Reify(ctx, c)
		assert.Nil(t, err)
		assert.NotNil(t, fromPid)

		select {
		case msg := <-wg:
			a, ok := msg.Annotation(timer.AnnotationTimerFiredCount)
			assert.True(t, ok)
			assert.Equal(t, "1", a)
		case <-time.After(5 * time.Second):
			assert.Fail(t, "timeout")
		}
	})
}

const simpleYAML = `
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
	support.Run(t, "run", func(t *testing.T, ctx context.Context, c camel.Context) {
		t.Helper()

		content := uuid.New()
		wg := make(chan camel.Message)

		c.Registry().Set("consumer-1", func(message camel.Message) {
			message.SetContent(content)
		})
		c.Registry().Set("consumer-2", func(message camel.Message) {
			wg <- message
		})

		err := c.LoadRoutes(ctx, strings.NewReader(simpleYAML))
		assert.Nil(t, err)

		select {
		case msg := <-wg:
			a, ok := msg.Annotation(timer.AnnotationTimerFiredCount)
			assert.True(t, ok)
			assert.Equal(t, "1", a)
			assert.Equal(t, content, msg.Content())
		case <-time.After(5 * time.Second):
			assert.Fail(t, "timeout")
		}
	})
}

const simpleWASM = `
- route:
    from:
      uri: "timer:foo"
      steps:
        - process:
            ref: "consumer-1"
        - transform:
            wasm: 
              path: "../../etc/fn/simple_process.wasm"
        - process:
            ref: "consumer-2"
`

func TestSimpleWASM(t *testing.T) {
	support.Run(t, "run", func(t *testing.T, ctx context.Context, c camel.Context) {
		t.Helper()

		wg := make(chan camel.Message)

		c.Registry().Set("consumer-1", func(message camel.Message) {
			_ = message.SetSubject("consumer-1")
		})
		c.Registry().Set("consumer-2", func(message camel.Message) {
			wg <- message
		})

		err := c.LoadRoutes(ctx, strings.NewReader(simpleWASM))
		assert.Nil(t, err)

		select {
		case msg := <-wg:
			assert.Equal(t, "consumer-1", msg.GetSubject())

			c, ok := msg.Content().([]byte)
			assert.True(t, ok)
			assert.Equal(t, "hello from wasm", string(c))

		case <-time.After(5 * time.Second):
			assert.Fail(t, "timeout")
		}
	})
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
`

func TestSimpleKafka(t *testing.T) {
	support.Run(t, "run", func(t *testing.T, ctx context.Context, c camel.Context) {
		t.Helper()

		content := uuid.New()

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

		ac, err := container.Admin(ctx)
		assert.Nil(t, err)

		tp, err := ac.CreateTopic(ctx, 3, 1, nil, "foo")
		assert.Nil(t, err)
		assert.Nil(t, tp.Err)

		c.Registry().Set("consumer-1", func(message camel.Message) {
			message.SetContent(content)
		})

		err = c.LoadRoutes(ctx, strings.NewReader(simpleKafka))
		assert.Nil(t, err)

		RegisterTestingT(t)

		Eventually(func(g Gomega) {
			f := cl.PollFetches(ctx)

			Expect(f.Errors()).To(BeEmpty())
			Expect(f.NumRecords()).To(Equal(1))
			Expect(string(f.Records()[0].Value)).To(Equal(content))
		}).Should(Succeed())
	})
}

const simpleKafkaWASM = `
- route:
    from:
      uri: "timer:foo"
      steps:
        - transform:
            wasm: 
              path: "../../etc/fn/simple_process.wasm"
        - to:
            uri: "kafka:foo"
            parameters:
              brokers: "localhost:9092"
`

func TestSimpleKafkaWASM(t *testing.T) {
	support.Run(t, "run", func(t *testing.T, ctx context.Context, c camel.Context) {
		t.Helper()

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

		ac, err := container.Admin(ctx)
		assert.Nil(t, err)

		tp, err := ac.CreateTopic(ctx, 3, 1, nil, "foo")
		assert.Nil(t, err)
		assert.Nil(t, tp.Err)

		err = c.LoadRoutes(ctx, strings.NewReader(simpleKafkaWASM))
		assert.Nil(t, err)

		RegisterTestingT(t)

		Eventually(func(g Gomega) {
			f := cl.PollFetches(ctx)

			Expect(f.Errors()).To(BeEmpty())
			Expect(f.NumRecords()).To(Equal(1))
			Expect(string(f.Records()[0].Value)).To(Equal("hello from wasm"))
		}).Should(Succeed())
	})
}

const simpleComponentWASM = `
- route:
    from:
      uri: "timer:foo"
      steps:
        - process:
            ref: "consumer-1"
        - to:
            uri: "wasm:../../etc/fn/simple_logger.wasm"
        - process:
            ref: "consumer-2"
`

func TestSimpleComponentWASM(t *testing.T) {
	support.Run(t, "run", func(t *testing.T, ctx context.Context, c camel.Context) {
		t.Helper()

		wg := make(chan camel.Message)

		c.Registry().Set("consumer-1", func(message camel.Message) {
			message.SetContent("consumer-1")
		})
		c.Registry().Set("consumer-2", func(message camel.Message) {
			wg <- message
		})

		err := c.LoadRoutes(ctx, strings.NewReader(simpleComponentWASM))
		assert.Nil(t, err)

		select {
		case msg := <-wg:
			c, ok := msg.Content().(string)
			assert.True(t, ok)
			assert.Equal(t, "consumer-1", c)

		case <-time.After(5 * time.Second):
			assert.Fail(t, "timeout")
		}
	})
}

const simpleComponentImageWASM = `
- route:
    from:
      uri: "timer:foo"
      steps:
        - process:
            ref: "consumer-1"
        - to:
            uri: "wasm:etc/fn/simple_logger.wasm?image=docker.io/lburgazzoli/camel-go:latest"
        - process:
            ref: "consumer-2"
`

func TestSimpleComponentImageWASM(t *testing.T) {
	support.Run(t, "run", func(t *testing.T, ctx context.Context, c camel.Context) {
		t.Helper()

		wg := make(chan camel.Message)

		c.Registry().Set("consumer-1", func(message camel.Message) {
			message.SetContent("consumer-1")
		})
		c.Registry().Set("consumer-2", func(message camel.Message) {
			wg <- message
		})

		err := c.LoadRoutes(ctx, strings.NewReader(simpleComponentImageWASM))
		assert.Nil(t, err)

		select {
		case msg := <-wg:
			c, ok := msg.Content().(string)
			assert.True(t, ok)
			assert.Equal(t, "consumer-1", c)

		case <-time.After(5 * time.Second):
			assert.Fail(t, "timeout")
		}
	})
}

const simpleMQTT = `
- route:
    from:
      uri: "timer:foo"
      steps:
        - process:
            ref: "consumer-1"
        - to:
            uri: "mqtt:foo"
            parameters:
              broker: "localhost:9092"
`

func TestSimpleMQTT(t *testing.T) {
	support.Run(t, "run", func(t *testing.T, ctx context.Context, c camel.Context) {
		t.Helper()
		t.Skipf("TBD")

		content := uuid.New()

		c.Registry().Set("consumer-1", func(message camel.Message) {
			message.SetContent(content)
		})

		err := c.LoadRoutes(ctx, strings.NewReader(simpleMQTT))
		assert.Nil(t, err)

		RegisterTestingT(t)

		Eventually(func(g Gomega) {
			// f := cl.PollFetches(ctx)

			// Expect(f.Errors()).To(BeEmpty())
			// Expect(f.NumRecords()).To(Equal(1))
			// Expect(string(f.Records()[0].Value)).To(Equal(content))
		}).Should(Succeed())
	})
}
