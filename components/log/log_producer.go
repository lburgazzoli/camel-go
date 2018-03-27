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

func newLogProducer(endpoint *logEndpoint, logger *zerolog.Logger) *logProducer {
	p := logProducer{
		endpoint:  endpoint,
		logger:    logger,
		processor: camel.NewProcessorSource(),
	}

	return &p
}

type logProducer struct {
	endpoint  *logEndpoint
	processor camel.Processor
	logger    *zerolog.Logger
}

func (producer *logProducer) Start() {
	producer.processor.Subscribe(producer.process)
}

func (producer *logProducer) Stop() {
}

func (producer *logProducer) Endpoint() camel.Endpoint {
	return producer.endpoint
}

func (producer *logProducer) Processor() camel.Processor {
	return producer.processor
}

func (producer *logProducer) process(exchange *camel.Exchange) {
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
