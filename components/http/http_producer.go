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
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	ghttp "net/http"
	gurl "net/url"
	"strings"

	"github.com/lburgazzoli/camel-go/camel"
	"github.com/lburgazzoli/camel-go/logger"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/processor"

	"github.com/rs/zerolog"
)

// ==========================
//
// Producer
//
// ==========================

func newHTTPProducer(endpoint *httpEndpoint) *httpProducer {
	p := httpProducer{
		logger:    logger.New("http.Producer"),
		endpoint:  endpoint,
		transport: endpoint.transport,
		client:    endpoint.client,
		converter: endpoint.component.context.TypeConverter(),
	}

	p.processor = processor.NewProcessingPipeline(p.process)

	return &p
}

type httpProducer struct {
	logger    zerolog.Logger
	endpoint  *httpEndpoint
	converter api.TypeConverter
	processor api.Processor
	transport *ghttp.Transport
	client    *ghttp.Client
}

func (producer *httpProducer) Start() {
	if producer.transport == nil {
		producer.transport = &ghttp.Transport{
			Dial: (&net.Dialer{
				Timeout: producer.endpoint.connectionTimeout,
			}).Dial,
		}
	}

	if producer.client == nil {
		producer.client = &ghttp.Client{
			Timeout:   producer.endpoint.requestTimeout,
			Transport: producer.transport,
		}
	}
}

func (producer *httpProducer) Stop() {
}

func (producer *httpProducer) Stage() api.ServiceStage {
	return api.ServiceStageProducer
}

func (producer *httpProducer) Endpoint() api.Endpoint {
	return producer.endpoint
}

func (producer *httpProducer) Processor() api.Processor {
	return producer.processor
}

//TODO: error handling
func (producer *httpProducer) process(exchange api.Exchange) {

	// compute the url
	url := fmt.Sprintf("%s://%s:%d%s", producer.endpoint.scheme, producer.endpoint.host, producer.endpoint.port, producer.endpoint.path)

	// Check if a method is set on the endpoint or on the exchange header
	if producer.endpoint.method == "" {
		headerMethod, found := exchange.Headers().Lookup(HTTP_METHOD)
		if found {
			// Check if method is valid
			if headerMethod != ghttp.MethodGet &&
				headerMethod != ghttp.MethodPost &&
				headerMethod != ghttp.MethodPut &&
				headerMethod != ghttp.MethodDelete &&
				headerMethod != ghttp.MethodOptions &&
				headerMethod != ghttp.MethodConnect &&
				headerMethod != ghttp.MethodHead &&
				headerMethod != ghttp.MethodPatch &&
				headerMethod != ghttp.MethodTrace {
				// do nothing here for the moment, we should fail the exchange
				producer.logger.Error().Msg("invalid HTTP method")
				return
			}
			producer.endpoint.method = headerMethod.(string)
		} else {
			// Default to GET
			producer.endpoint.method = "GET"
		}
	}

	// Convert the body to a byte array
	body := []byte(fmt.Sprintf("%+v", exchange.Body()))
	if bodyBytes, ok := exchange.Body().([]byte); ok {
		body = bodyBytes
	}

	// Create the request
	req, err := ghttp.NewRequest(producer.endpoint.method, url, bytes.NewBuffer(body))
	if err != nil {
		// do nothing here for the moment, we should fail the exchange
		producer.logger.Error().Msg(err.Error())
	} else {

		// Set the query parameters
		if query, found := exchange.Headers().Lookup(HTTP_QUERY); found {
			if queryVal, ok := query.(gurl.Values); ok {
				req.URL.RawQuery = queryVal.Encode()
			} else {
				if strQuery, ok := query.(string); ok {
					// Parse the query string
					if queryVal, err := gurl.ParseQuery(strQuery); err == nil {
						req.URL.RawQuery = queryVal.Encode()
					}
				}
			}
		}

		// Set the headers
		exchange.Headers().ForEach(func(key string, val interface{}) {
			// Skip Camel internal headers
			if !strings.HasPrefix(key, camel.CAMEL_HEADER) {

				if values, ok := val.([]string); ok {
					for _, v := range values {
						req.Header.Add(key, v)
					}
				} else {
					if v, err := producer.converter(val, camel.TypeString); err == nil {
						req.Header.Set(key, v.(string))
					}
				}
			}
		})

		response, err := producer.client.Do(req)

		if err != nil {
			// do nothing here for the moment, we should fail the exchange
			producer.logger.Error().Msg(err.Error())
		}

		defer response.Body.Close()

		exchange.Headers().Bind(HTTP_STATUS_CODE, fmt.Sprint(response.StatusCode))
		exchange.Headers().Bind(HTTP_CONTENT_LENGTH, response.ContentLength)

		for k, v := range response.Header {
			if len(v) >= 1 {
				exchange.Headers().Bind(k, v)
			}
		}

		// we should handle status code, set headers & so on here.
		if response.StatusCode == ghttp.StatusOK {
			bytes, _ := ioutil.ReadAll(response.Body)
			exchange.SetBody(string(bytes))
		}
	}
}
