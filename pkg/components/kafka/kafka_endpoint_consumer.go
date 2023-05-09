////go:build components_kafka || components_all

package kafka

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Consumer struct {
	components.DefaultConsumer

	endpoint *Endpoint
	client   *kgo.Client
	running  atomic.Bool
}

func (c *Consumer) Endpoint() camel.Endpoint {
	return c.endpoint
}

func (c *Consumer) Start(ctx context.Context) error {

	if c.running.CompareAndSwap(false, true) {
		cl, err := c.endpoint.newClient()
		if err != nil {
			return err
		}

		c.client = cl

		go func() {
			if !c.running.Load() {
				return
			}

			c.poll(ctx)
		}()
	}

	return nil
}

func (c *Consumer) Stop(context.Context) error {
	return nil
}

func (c *Consumer) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		_ = c.Start(context.Background())
	case *actor.Stopping:
		_ = c.Stop(context.Background())
	case camel.Message:
		// ignore message,
		// TODO: may be used for transactions
		break
	}
}

func (c *Consumer) poll(ctx context.Context) {

	component := c.endpoint.Component()
	context := component.Context()

	for {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			// All errors are retried internally when fetching, but non-retriable errors are
			// returned from polls so that users can notice and take action.
			panic(fmt.Sprint(errs))
		}

		for it := fetches.RecordIter(); !it.Done(); {
			record := it.Next()

			m, err := message.New()
			if err != nil {
				panic(err)
			}

			_ = m.SetType("camel.kafka.record.received")
			_ = m.SetSource(component.Scheme())

			if len(record.Key) > 0 {
				_ = m.SetExtension("partitionkey", record.Key)

				// TODO: use type converter
				_ = m.SetSubject(string(record.Key))
			}

			for _, h := range record.Headers {
				// TODO: use type converter
				_ = m.SetExtension(h.Key, string(h.Value))
			}

			m.SetContent(record.Value)

			m.SetAnnotation(AnnotationOffset, strconv.FormatInt(record.Offset, 10))
			m.SetAnnotation(AnnotationPartition, strconv.FormatInt(int64(record.Partition), 10))
			m.SetAnnotation(AnnotationTopic, record.Topic)

			// TODO: use type converter
			m.SetAnnotation(AnnotationKey, string(record.Key))

			if err := context.SendTo(c.Target(), m); err != nil {
				panic(err)
			}
		}
	}
}
