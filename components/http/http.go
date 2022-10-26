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

import "github.com/lburgazzoli/camel-go/camel"

// ==========================
//
// Init
//
// ==========================

func init() {
	camel.RootContext.Registry().Bind("http", NewComponent())
}

// ==========================
//
// this is where constant
// should be set
//
// ==========================

// HTTP Method Header --
const HTTP_METHOD = "camel.http.method"

// HTTP Query Header --
const HTTP_QUERY = "camel.http.query"

// HTTP Request Path Header --
const HTTP_REQUEST_PATH = "camel.http.requestPath"

// HTTP Status Code Header --
const HTTP_STATUS_CODE = "camel.http.statusCode"

// HTTP Content Length Header --
const HTTP_CONTENT_LENGTH = "camel.http.contentLength"
