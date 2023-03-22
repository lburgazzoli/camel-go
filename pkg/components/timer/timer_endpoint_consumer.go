////go:build components_timer || components_all

package timer

import (
	"context"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/procyon-projects/chrono"
)

type Consumer struct {
	api.WithOutputs

	id        string
	endpoint  *Endpoint
	scheduler chrono.TaskScheduler
	task      chrono.ScheduledTask
	counter   uint64
	started   time.Time
}

func (c *Consumer) Endpoint() api.Endpoint {
	return c.endpoint
}

func (c *Consumer) ID() string {
	return c.id
}

func (c *Consumer) Start(context.Context) error {
	c.counter = 0
	c.started = time.Now()
	c.scheduler = chrono.NewDefaultTaskScheduler()

	t, err := c.scheduler.ScheduleWithFixedDelay(c.run, c.endpoint.config.Interval)
	if err != nil {
		return err
	}

	c.task = t

	return nil
}

func (c *Consumer) Stop(context.Context) error {
	if c.task != nil {
		c.task.Cancel()
		c.task = nil
	}
	if c.scheduler != nil {
		c.scheduler.Shutdown()
		c.scheduler = nil
	}

	c.counter = 0
	c.started = time.UnixMilli(0)

	return nil
}

func (c *Consumer) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		_ = c.Start(context.Background())
	case *actor.Stopping:
		_ = c.Stop(context.Background())
	}
}

func (c *Consumer) run(_ context.Context) {
	component := c.endpoint.Component()
	context := component.Context()

	for _, o := range c.Outputs() {
		m, err := message.New()
		if err != nil {
			panic(err)
		}

		_ = m.SetType("camel.timer.triggered")
		_ = m.SetSource(component.Scheme())

		m.SetAnnotation(AnnotationTimerFiredCount, strconv.FormatUint(atomic.AddUint64(&c.counter, 1), 10))
		m.SetAnnotation(AnnotationTimerStarted, strconv.FormatInt(c.started.UnixMilli(), 19))
		m.SetAnnotation(AnnotationTimerName, c.endpoint.config.Remaining)

		if err := context.SendTo(o, m); err != nil {
			panic(err)
		}
	}
}
