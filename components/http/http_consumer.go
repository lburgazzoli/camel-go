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
	"io/ioutil"
	ghttp "net/http"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/camel"
	"github.com/lburgazzoli/camel-go/processor"
)

// ==========================
//
// Producer
//
// ==========================

func newHTTPConsumer(endpoint *httpEndpoint) *httpConsumer {
	c := httpConsumer{
		endpoint: endpoint,
		// TODO: this is ugly
		processor: processor.NewProcessingPipeline(func(api.Exchange) {
		}),
	}

	return &c
}

type httpConsumer struct {
	endpoint  *httpEndpoint
	processor api.Processor
	server    *ghttp.Server
}

func (consumer *httpConsumer) Start() {
	if consumer.server != nil {
		return
	}

	mux := ghttp.NewServeMux()
	mux.HandleFunc(consumer.endpoint.path, func(w ghttp.ResponseWriter, r *ghttp.Request) {
		w.WriteHeader(ghttp.StatusOK)

		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		exchange := camel.NewExchange(consumer.endpoint.component.context)
		exchange.SetBody(string(body))

		consumer.processor.Publish(exchange)
	})

	consumer.server = &ghttp.Server{Addr: "", Handler: mux}
	go func() {
		consumer.server.ListenAndServe()
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
