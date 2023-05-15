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
	camelCtx := component.Context()

	for {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			// All errors are retried internally when fetching, but non-retriable errors are
			// returned from polls so that users can notice and take action.
			panic(fmt.Sprint(errs))
		}

		for it := fetches.RecordIter(); !it.Done(); {
			record := it.Next()

			m := camelCtx.NewMessage()

			m.SetType("camel.kafka.record.received")
			m.SetSource(component.Scheme())

			if len(record.Key) > 0 {
				m.SetHeader("partitionkey", record.Key)

				// TODO: use type converter
				m.SetSubject(string(record.Key))
			}

			for _, h := range record.Headers {
				// TODO: use type converter
				m.SetHeader(h.Key, string(h.Value))
			}

			m.SetContent(record.Value)

			_ = m.SetAttribute(AttributeOffset, strconv.FormatInt(record.Offset, 10))
			_ = m.SetAttribute(AttributePartition, strconv.FormatInt(int64(record.Partition), 10))
			_ = m.SetAttribute(AttributeTopic, record.Topic)
			_ = m.SetAttribute(AttributeKey, string(record.Key))

			if err := camelCtx.SendTo(c.Target(), m); err != nil {
				panic(err)
			}
		}
	}
}
