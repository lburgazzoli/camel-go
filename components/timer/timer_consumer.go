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

type timerConsumer struct {
	endpoint *timerEndpoint
	pipe     *camel.Pipe
	ticker   *time.Ticker
}

// Endpoint --
func (consumer *timerConsumer) Endpoint() camel.Endpoint {
	return consumer.endpoint
}

// Start --
func (consumer *timerConsumer) Start() {
	consumer.ticker = time.NewTicker(consumer.endpoint.period)
	go func() {
		var counter uint64

		for t := range consumer.ticker.C {
			exchange := camel.NewExchange(consumer.endpoint.component.context)

			counter++

			exchange.SetHeader("timer.fire.time", t.UTC())
			exchange.SetHeader("timer.fire.count", counter)
			exchange.SetBody(nil)

			consumer.pipe.Publish(exchange)
		}
	}()
}

// Stop
func (consumer *timerConsumer) Stop() {
	if consumer.ticker != nil {
		consumer.ticker.Stop()
	}
}
