////go:build components_timer || components_all

package timer

import (
	"context"
	"sync/atomic"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/procyon-projects/chrono"
)

type Consumer struct {
	api.WithOutputs

	endpoint  *Endpoint
	scheduler chrono.TaskScheduler
	task      chrono.ScheduledTask
	counter   uint64
}

func (c *Consumer) Endpoint() api.Endpoint {
	return c.endpoint
}

func (c *Consumer) Start() error {
	c.counter = 0
	c.scheduler = chrono.NewDefaultTaskScheduler()

	t, err := c.scheduler.ScheduleWithFixedDelay(c.run, c.endpoint.config.Interval)
	if err != nil {
		return err
	}

	c.task = t

	return nil
}

func (c *Consumer) Stop() error {
	if c.task != nil {
		c.task.Cancel()
	}
	if c.scheduler != nil {
		c.scheduler.Shutdown()
	}

	return nil
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

		m.SetAnnotation("counter", atomic.AddUint64(&c.counter, 1))

		context.Send(o, m)
	}
}
