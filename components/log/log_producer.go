package log

import (
	"github.com/lburgazzoli/camel-go/api"
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
		endpoint: endpoint,
		logger:   logger,
	}

	p.processor = api.NewProcessingPipeline(p.process)

	return &p
}

type logProducer struct {
	endpoint  *logEndpoint
	processor api.Processor
	converter api.TypeConverter
	logger    *zerolog.Logger
}

func (producer *logProducer) Start() {
}

func (producer *logProducer) Stop() {
}

func (producer *logProducer) Endpoint() api.Endpoint {
	return producer.endpoint
}

func (producer *logProducer) Processor() api.Processor {
	return producer.processor
}

func (producer *logProducer) process(exchange api.Exchange) {
	lg := producer.logger.WithLevel(producer.endpoint.level)
	tc := producer.endpoint.component.context.TypeConverter()
	body := exchange.Body()

	if producer.endpoint.logHeaders {
		d := zerolog.Dict()

		exchange.Headers().Range(func(key string, val interface{}) bool {
			if str, err := tc(val, camel.TypeString); str != nil && err != nil {
				d.Str(key, str.(string))
			} else {
				d.Interface(key, val)
			}

			return true
		})

		if str, err := tc(body, camel.TypeString); str != nil && err != nil && body != nil {
			lg.Dict("headers", d).Msgf("body: %s", str)
		} else {
			lg.Dict("headers", d).Msgf("body: %+v", body)
		}
	} else {
		if str, err := tc(body, camel.TypeString); str != nil && err != nil && body != nil {
			lg.Msgf("body: %s", body)
		} else {
			lg.Msgf("body: %+v", body)
		}
	}
}
