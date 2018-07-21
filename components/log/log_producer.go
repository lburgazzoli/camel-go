// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/camel"
	"github.com/lburgazzoli/camel-go/processor"
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

	p.processor = processor.NewProcessingPipeline(p.process)

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

func (producer *logProducer) Stage() api.ServiceStage {
	return api.ServiceStageProducer
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
			lg.Dict("headers", d).Msgf("%s", str)
		} else {
			lg.Dict("headers", d).Msgf("%+v", body)
		}
	} else {
		if str, err := tc(body, camel.TypeString); str != nil && err != nil && body != nil {
			lg.Msgf("%s", body)
		} else {
			lg.Msgf("%+v", body)
		}
	}
}
