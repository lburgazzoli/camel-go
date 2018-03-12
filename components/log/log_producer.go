package log

import (
	"github.com/lburgazzoli/camel-go/camel"
	"github.com/rs/zerolog"
)

// ==========================
//
// Producer
//
// ==========================

type logProducer struct {
	endpoint *logEndpoint
	logger   *zerolog.Logger
}

func (producer *logProducer) Endpoint() camel.Endpoint {
	return producer.endpoint
}

func (endpoint *logProducer) Start() {
}

func (endpoint *logProducer) Stop() {
}

func (producer *logProducer) Process(exchange *camel.Exchange) {
	if producer.endpoint.logHeaders {
		l := producer.logger.WithLevel(producer.endpoint.level)
		d := zerolog.Dict()

		for k, v := range exchange.Headers() {
			d.Interface(k, v)
		}

		l.Dict("headers", d).Msgf("%+v", exchange.Body())
	} else {
		producer.logger.WithLevel(producer.endpoint.level).Msgf("%+v", exchange.Body())
	}
}
