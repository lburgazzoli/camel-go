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
	ghttp "net/http"
	"time"

	"github.com/lburgazzoli/camel-go/api"
)

// ==========================
//
// Endpoint
//
// ==========================

type httpEndpoint struct {
	component         *Component
	scheme            string
	host              string
	port              int
	path              string
	method            string
	connectionTimeout time.Duration
	requestTimeout    time.Duration
	transport         *ghttp.Transport
	client            *ghttp.Client
}

func (endpoint *httpEndpoint) Start() {
}

func (endpoint *httpEndpoint) Stop() {
}

func (endpoint *httpEndpoint) Stage() api.ServiceStage {
	return api.ServiceStageEndpoint
}

func (endpoint *httpEndpoint) Component() api.Component {
	return endpoint.component
}

func (endpoint *httpEndpoint) CreateProducer() (api.Producer, error) {
	return newHTTPProducer(endpoint), nil
}

func (endpoint *httpEndpoint) CreateConsumer() (api.Consumer, error) {
	return newHTTPConsumer(endpoint), nil
}

// ==========================
//
// Options
//
// ==========================

func (endpoint *httpEndpoint) SetMethod(method string) {
	endpoint.method = method
}

func (endpoint *httpEndpoint) SetConnectionTimeout(timeout time.Duration) {
	endpoint.connectionTimeout = timeout
}

func (endpoint *httpEndpoint) SetRequestTimeout(timeout time.Duration) {
	endpoint.requestTimeout = timeout
}

func (endpoint *httpEndpoint) SetTransport(transport *ghttp.Transport) {
	endpoint.transport = transport
}

func (endpoint *httpEndpoint) SetClient(client *ghttp.Client) {
	endpoint.client = client
}
