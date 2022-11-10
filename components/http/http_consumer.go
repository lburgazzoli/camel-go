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

package http

import (
	"fmt"
	"io/ioutil"
	ghttp "net/http"
	"strings"

	"github.com/rs/zerolog"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/camel"
	"github.com/lburgazzoli/camel-go/logger"
	"github.com/lburgazzoli/camel-go/processor"
)

// ==========================
//
// Producer
//
// ==========================

func newHTTPConsumer(endpoint *httpEndpoint) *httpConsumer {
	c := httpConsumer{
		logger:   logger.New("http.Consumer"),
		endpoint: endpoint,
		// TODO: this is ugly
		processor: processor.NewProcessingPipeline(func(api.Exchange) {
		}),
	}

	return &c
}

type httpConsumer struct {
	logger    zerolog.Logger
	endpoint  *httpEndpoint
	processor api.Processor
	server    *ghttp.Server
}

func (consumer *httpConsumer) Start() {
	if consumer.server != nil {
		return
	}

	// compute the url
	url := fmt.Sprintf("%s:%d", consumer.endpoint.host, consumer.endpoint.port)

	mux := ghttp.NewServeMux()
	mux.HandleFunc(consumer.endpoint.path, func(w ghttp.ResponseWriter, r *ghttp.Request) {

		// check method
		if consumer.endpoint.method != "" && r.Method != consumer.endpoint.method {
			w.WriteHeader(ghttp.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		exchange := camel.NewExchange(consumer.endpoint.component.context)
		exchange.Headers().Bind(HTTP_REQUEST_PATH, r.URL.Path)
		exchange.Headers().Bind(HTTP_METHOD, r.Method)
		exchange.Headers().Bind(HTTP_QUERY, r.URL.RawQuery)

		for queryKey, queryValues := range r.URL.Query() {
			if len(queryValues) == 1 {
				exchange.Headers().Bind(queryKey, queryValues[0])
			} else {
				exchange.Headers().Bind(queryKey, queryValues)
			}
		}

		for headerKey, headerValue := range r.Header {
			if !strings.HasPrefix(headerKey, camel.CAMEL_HEADER) {
				exchange.Headers().Bind(headerKey, headerValue)
			}
		}

		exchange.SetBody(string(body))

		consumer.processor.Publish(exchange)

		// Wait for the returned exchange
		ch := make(chan api.Exchange)
		subscription := consumer.processor.SubscribeReturn(func(retExchange api.Exchange) {
			ch <- retExchange
		})
		retExchange := <-ch
		subscription.Cancel()

		// Write the headers
		retExchange.Headers().ForEach(func(key string, value any) {

			switch {
			case strings.HasPrefix(key, camel.CAMEL_HEADER):
				// Skip Camel internal headers

			case key == "Content-Length":
				// Skip Content-Length header

			default:
				camel.ForEachIn(value, func(_, v any) {
					w.Header().Add(key, fmt.Sprintf("%v", v))
				})
			}
		})

		if retExchange.IsFailed() {
			w.WriteHeader(ghttp.StatusInternalServerError)
			// Write the http body with the error message
			if _, err := w.Write([]byte(fmt.Sprintf("%+v", retExchange.Error()))); err != nil {
				consumer.logger.Error().Err(err).Msg("Error writing response")
			}
		} else {
			w.WriteHeader(ghttp.StatusOK)
			// Write the http body with the message body
			if _, err := w.Write([]byte(fmt.Sprintf("%+v", retExchange.Body()))); err != nil {
				consumer.logger.Error().Err(err).Msg("Error writing response")
			}
		}

	})

	consumer.server = &ghttp.Server{Addr: url, Handler: mux}

	go func() {
		consumer.logger.Debug().Msgf("Start listening %+v requests on address: %+v", consumer.endpoint.method, consumer.server.Addr)
		if err := consumer.server.ListenAndServe(); err != nil {
			consumer.logger.Fatal().Msg(err.Error())
		}
		consumer.logger.Debug().Msgf("Stop serving %s", consumer.server.Addr)
	}()
}

func (consumer *httpConsumer) Stop() {
	if consumer.server != nil {
		consumer.server.Close()
	}
}

func (consumer *httpConsumer) Stage() api.ServiceStage {
	return api.ServiceStageConsumer
}

func (consumer *httpConsumer) Endpoint() api.Endpoint {
	return consumer.endpoint
}

func (consumer *httpConsumer) Processor() api.Processor {
	return consumer.processor
}
