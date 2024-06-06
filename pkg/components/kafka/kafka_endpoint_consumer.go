// //go:build components_kafka || components_all

package kafka

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/lburgazzoli/camel-go/pkg/cloudevents"

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
		cl, err := c.endpoint.newClient(
			kgo.BlockRebalanceOnPoll(),
		)

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
	if c.running.CompareAndSwap(true, false) {
		c.client.Close()
		c.client = nil
	}

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
		// TODO: may be used for ack
		break
	}
}

func (c *Consumer) poll(ctx context.Context) {
	for {
		closed, err := c.pollRecords(ctx)
		if closed {
			return
		}

		if err != nil {
			panic(err)
		}
	}
}

func (c *Consumer) pollRecords(ctx context.Context) (bool, error) {
	defer c.client.AllowRebalance()

	fetches := c.client.PollRecords(ctx, DefaultPollSize)
	if fetches.IsClientClosed() {
		return false, nil
	}

	if fe := fetches.Errors(); len(fe) > 0 {
		// All errors are retried internally when fetching, but non-retriable errors are
		// returned from polls so that users can notice and take action.
		allErrors := make([]error, len(fe))
		for i := range fe {
			allErrors[i] = fmt.Errorf(
				"non-retriable error polling for records on topic: %s, partition: %d, error: %w",
				fe[i].Topic,
				fe[i].Partition,
				fe[i].Err)
		}

		return true, errors.Join(allErrors...)
	}

	component := c.endpoint.Component()
	camelCtx := component.Context()

	for it := fetches.RecordIter(); !it.Done(); {
		record := it.Next()

		m := camelCtx.NewMessage()

		m.SetType("camel.kafka.record.received")
		m.SetSource(component.Scheme())

		m.SetHeader(cloudevents.ExtensionSequence, fmt.Sprintf("%d-%d", record.Partition, record.Offset))

		if len(record.Key) > 0 {
			m.SetHeader(cloudevents.ExtensionPartitionKey, string(record.Key))
			m.SetSubject(string(record.Key))
		}

		for _, h := range record.Headers {
			m.SetHeader(h.Key, string(h.Value))
		}

		m.SetContent(record.Value)

		m.SetAttribute(AttributeOffset, strconv.FormatInt(record.Offset, 10))
		m.SetAttribute(AttributePartition, strconv.FormatInt(int64(record.Partition), 10))
		m.SetAttribute(AttributeTopic, record.Topic)
		m.SetAttribute(AttributeKey, string(record.Key))

		if err := camelCtx.SendTo(c.Target(), m); err != nil {
			return true, fmt.Errorf(
				"error sending to target with id: %s, address: %s, error: %w",
				c.Target().GetId(),
				c.Target().GetAddress(),
				err)
		}
	}

	return true, nil
}
