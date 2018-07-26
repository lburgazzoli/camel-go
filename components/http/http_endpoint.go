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
	"net/url"
	"strconv"
	"time"

	"github.com/lburgazzoli/camel-go/api"
)

// ==========================
//
// Options
//
// ==========================

// Options --
type EndpointOptions struct {
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

// SetMethod --
func (options *EndpointOptions) SetMethod(method string) {
	options.method = method
}

// SetConnectionTimeout --
func (options *EndpointOptions) SetConnectionTimeout(timeout time.Duration) {
	options.connectionTimeout = timeout
}

// SetRequestTimeout --
func (options *EndpointOptions) SetRequestTimeout(timeout time.Duration) {
	options.requestTimeout = timeout
}

// SetTransport --
func (options *EndpointOptions) SetTransport(transport *ghttp.Transport) {
	options.transport = transport
}

// SetClient --
func (options *EndpointOptions) SetClient(client *ghttp.Client) {
	options.client = client
}

// ==========================
//
// Functional Options
//
// ==========================

// EndpointOption --
type EndpointOption func(*EndpointOptions)

// Method --
func Method(value string) EndpointOption {
	return func(args *EndpointOptions) {
		args.method = value
	}
}

// ConnectionTimeout --
func ConnectionTimeout(value time.Duration) EndpointOption {
	return func(args *EndpointOptions) {
		args.connectionTimeout = value
	}
}

// RequestTimeout --
func RequestTimeout(value time.Duration) EndpointOption {
	return func(args *EndpointOptions) {
		args.requestTimeout = value
	}
}

// Transport --
func Transport(value *ghttp.Transport) EndpointOption {
	return func(args *EndpointOptions) {
		args.transport = value
	}
}

// Client --
func Client(value *ghttp.Client) EndpointOption {
	return func(args *EndpointOptions) {
		args.client = value
	}
}

// ==========================
//
// Endpoint
//
// ==========================

func newEndpoint(component *Component, url url.URL, setters ...EndpointOption) (*httpEndpoint, error) {
	endpoint := httpEndpoint{}
	endpoint.component = component
	endpoint.method = ghttp.MethodGet
	endpoint.connectionTimeout = 10 * time.Second
	endpoint.requestTimeout = 60 * time.Second
	endpoint.path = url.Path
	endpoint.port = 80

	if url.Hostname() != "" {
		endpoint.host = url.Hostname()
	}

	if url.Port() != "" {
		port, err := strconv.Atoi(url.Port())
		if err != nil {
			return nil, err
		}

		endpoint.port = port
	}

	if endpoint.port == 443 && endpoint.scheme == "" {
		endpoint.scheme = "https"
	}

	if endpoint.scheme == "" {
		endpoint.scheme = "http"
	}

	// Apply options
	for _, setter := range setters {
		setter(&endpoint.EndpointOptions)
	}

	return &endpoint, nil
}

type httpEndpoint struct {
	EndpointOptions

	component *Component
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
