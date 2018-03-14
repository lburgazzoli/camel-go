package log

import (
	"github.com/lburgazzoli/camel-go/camel"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ==========================
//
// Producer
//
// ==========================

func newLogProducer(endpoint *logEndpoint, logger *zerolog.Logger) *logProducer {
	p := logProducer{
		endpoint: endpoint,
		logger:   logger,
		pipe:     &camel.Pipe{},
	}

	return &p
}

type logProducer struct {
	endpoint *logEndpoint
	pipe     *camel.Pipe
	logger   *zerolog.Logger
}

func (producer *logProducer) Start() {
	go func() {
		for {
			select {
			case exchange, ok := <-producer.pipe.In:
				if !ok {
					log.Warn().Msgf("Channel %+v is not ready", producer.pipe.In)
				} else {
					producer.process(exchange)
					producer.pipe.Publish(exchange)
				}
			case <-producer.pipe.Done:
				log.Info().Msg("done")
				return
			}
		}
	}()
}

func (producer *logProducer) Stop() {
}

func (producer *logProducer) Endpoint() camel.Endpoint {
	return producer.endpoint
}

func (producer *logProducer) Pipe() *camel.Pipe {
	return producer.pipe
}

func (producer *logProducer) process(exchange *camel.Exchange) *camel.Exchange {
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

	return exchange
}
