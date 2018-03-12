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
	endpoint  *timerEndpoint
	processor camel.Processor
	ticker    *time.Ticker
}

// Endpoint --
func (producer *timerConsumer) Endpoint() camel.Endpoint {
	return producer.endpoint
}

// Start --
func (producer *timerConsumer) Start() {
	producer.ticker = time.NewTicker(producer.endpoint.period)
	go func() {
		var counter uint64

		for t := range producer.ticker.C {
			exchange := camel.NewExchange(producer.endpoint.component.context)

			counter++

			exchange.SetHeader("timer.fire.time", t.UTC())
			exchange.SetHeader("timer.fire.count", counter)
			exchange.SetBody(nil)

			producer.processor.Process(exchange)
		}
	}()
}

// Stop
func (producer *timerConsumer) Stop() {
	if producer.ticker != nil {
		producer.ticker.Stop()
	}
}
