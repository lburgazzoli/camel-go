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
	"net"
	ghttp "net/http"
	"strings"

	"github.com/lburgazzoli/camel-go/camel"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/processor"

	zlog "github.com/rs/zerolog/log"
)

// ==========================
//
// Producer
//
// ==========================

func newHTTPProducer(endpoint *httpEndpoint) *httpProducer {
	p := httpProducer{
		endpoint:  endpoint,
		transport: endpoint.transport,
		client:    endpoint.client,
		converter: endpoint.component.context.TypeConverter(),
	}

	p.processor = processor.NewProcessingPipeline(p.process)

	return &p
}

type httpProducer struct {
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

	var url string

	// compute the url
	if producer.endpoint.path == "" {
		url = fmt.Sprintf("%s://%s:%d", producer.endpoint.scheme, producer.endpoint.host, producer.endpoint.port)
	} else {
		url = fmt.Sprintf("%s://%s:%d/%s", producer.endpoint.scheme, producer.endpoint.host, producer.endpoint.port, producer.endpoint.path)
	}

	req, err := ghttp.NewRequest(producer.endpoint.method, url, nil)
	if err != nil {
		// do nothing here for the moment, we should fail tyhe exchange
		zlog.Error().Msg(err.Error())
	} else {
		exchange.Headers().ForEach(func(key string, val interface{}) {
			if strings.HasPrefix(key, HTTPHeaderPrefix) && len(key) > HTTPHeaderPrefixLen {
				key = key[HTTPHeaderPrefixLen:]

				if v, err := producer.converter(val, camel.TypeString); err == nil {
					req.Header.Set(key, v.(string))
				}
			}
		})

		response, err := producer.client.Do(req)

		if err != nil {
			// do nothing here for the moment, we should fail tyhe exchange
			zlog.Error().Msg(err.Error())
		}

		defer response.Body.Close()

		exchange.Headers().Bind("http.StatusCode", response.StatusCode)
		exchange.Headers().Bind("http.ContentLength", response.ContentLength)

		for k, v := range response.Header {
			if len(v) >= 1 {
				exchange.Headers().Bind("http."+k, v[0])
			}
		}

		// we should handle status code, set headers & so on here.
		if response.StatusCode == ghttp.StatusOK {
			bytes, _ := ioutil.ReadAll(response.Body)
			exchange.SetBody(string(bytes))
		}
	}
}
