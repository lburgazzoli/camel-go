// //go:build components_kafka || components_all

package kafka

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/twmb/franz-go/pkg/kgo"

	camel "github.com/lburgazzoli/camel-go/pkg/api"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	. "github.com/onsi/gomega"

	// test support.
	"github.com/lburgazzoli/camel-go/pkg/util/tests/containers/kafka"

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
              brokers: "${kafka.broker}"
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

		require.NoError(t, container.Start(ctx))

		ac, err := container.Admin(ctx)
		require.NoError(t, err)

		tp, err := ac.CreateTopic(ctx, 3, 1, nil, "foo")
		require.NoError(t, err)
		require.NoError(t, tp.Err)

		cl, err := container.Client(
			ctx,
			kgo.ConsumeTopics("foo"),
			kgo.ConsumerGroup(uuid.New()),
		)

		require.NoError(t, err)

		defer cl.Close()

		props, err := container.Properties(ctx)
		require.NoError(t, err)

		c := camel.ExtractContext(ctx)

		err = c.Properties().Add(props)
		require.NoError(t, err)

		c.Registry().Set("consumer-1", func(_ context.Context, message camel.Message) error {
			message.SetContent(content)
			return nil
		})

		err = c.LoadRoutes(ctx, strings.NewReader(simpleKafka))
		require.NoError(t, err)

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
              brokers: "${kafka.broker}"
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

		require.NoError(t, container.Start(ctx))

		ac, err := container.Admin(ctx)
		require.NoError(t, err)

		tp, err := ac.CreateTopic(ctx, 3, 1, nil, "foo")
		require.NoError(t, err)
		require.NoError(t, tp.Err)

		cl, err := container.Client(
			ctx,
			kgo.ConsumeTopics("foo"),
			kgo.ConsumerGroup(uuid.New()),
		)

		require.NoError(t, err)

		defer cl.Close()

		props, err := container.Properties(ctx)
		require.NoError(t, err)

		c := camel.ExtractContext(ctx)

		err = c.Properties().Add(props)
		require.NoError(t, err)

		err = c.LoadRoutes(ctx, strings.NewReader(simpleKafkaWASM))
		require.NoError(t, err)

		RegisterTestingT(t)

		Eventually(func(g Gomega) {
			f := cl.PollFetches(ctx)

			Expect(f.Errors()).To(BeEmpty())
			Expect(f.NumRecords()).To(Equal(1))
			Expect(string(f.Records()[0].Value)).To(Equal("hello from wasm"))
		}).Should(Succeed())
	})
}
