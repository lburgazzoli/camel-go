//go:build components_timer || components_all

package timer

import (
	"context"
	"github.com/procyon-projects/chrono"
)

// Consumer ---
type Consumer struct {
	endpoint  *Endpoint
	scheduler chrono.TaskScheduler
	task      chrono.ScheduledTask
}

func (c *Consumer) Start() error {
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

func (c *Consumer) run(tx context.Context) {

}
