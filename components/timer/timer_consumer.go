package timer

import (
	"time"

	"github.com/lburgazzoli/camel-go/camel"
)

// ==========================
//
// Producer
//
// ==========================

func newTimerConsumer(endpoint *timerEndpoint) *timerConsumer {
	c := timerConsumer{
		endpoint:  endpoint,
		processor: camel.NewProcessorSource(),
	}

	return &c
}

type timerConsumer struct {
	endpoint  *timerEndpoint
	processor camel.Processor
	ticker    *time.Ticker
}

// Start --
func (consumer *timerConsumer) Start() {
	consumer.ticker = time.NewTicker(consumer.endpoint.period)
	go func() {
		var counter uint64

		for t := range consumer.ticker.C {
			exchange := camel.NewExchange(consumer.endpoint.component.context)

			counter++

			exchange.Headers().Bind("timer.fire.time", t.UTC())
			exchange.Headers().Bind("timer.fire.count", counter)
			exchange.SetBody(nil)

			consumer.processor.Publish(exchange)
		}
	}()
}

// Stop
func (consumer *timerConsumer) Stop() {
	if consumer.ticker != nil {
		consumer.ticker.Stop()
	}
}

// Endpoint --
func (consumer *timerConsumer) Endpoint() camel.Endpoint {
	return consumer.endpoint
}

func (consumer *timerConsumer) Processor() camel.Processor {
	return consumer.processor
}
