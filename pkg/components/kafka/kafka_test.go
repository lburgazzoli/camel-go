package kafka

import (
	"context"
	"strings"
	"testing"

	"github.com/twmb/franz-go/pkg/kgo"

	camel "github.com/lburgazzoli/camel-go/pkg/api"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/stretchr/testify/assert"

	// test support.
	"github.com/lburgazzoli/camel-go/pkg/util/tests/containers/kafka"

	// gomega.
	. "github.com/onsi/gomega"

	// enable components.
	_ "github.com/lburgazzoli/camel-go/pkg/components/log"
	_ "github.com/lburgazzoli/camel-go/pkg/components/timer"

	// enable processors.
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/process"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/to"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/transform"
)

const simpleKafka = `
- route:
    from:
      uri: "timer:foo"
      steps:
        - process:
            ref: "consumer-1"
        - to:
            uri: "log:info"
        - to:
            uri: "kafka:foo"
            parameters:
              brokers: "localhost:9092"
`

func TestSimpleKafka(t *testing.T) {
	support.Run(t, "run", func(t *testing.T, ctx context.Context) {
		t.Helper()

		content := uuid.New()

		container, err := kafka.NewContainer(ctx)
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

		c := camel.GetContext(ctx)

		c.Registry().Set("consumer-1", func(_ context.Context, message camel.Message) error {
			message.SetContent(content)
			return nil
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
              path: "../../../etc/wasm/fn/simple_process.wasm"
        - to:
            uri: "log:info"
        - to:
            uri: "kafka:foo"
            parameters:
              brokers: "localhost:9092"
`

func TestSimpleKafkaWASM(t *testing.T) {
	support.Run(t, "run", func(t *testing.T, ctx context.Context) {
		t.Helper()

		container, err := kafka.NewContainer(ctx)
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

		c := camel.GetContext(ctx)

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
