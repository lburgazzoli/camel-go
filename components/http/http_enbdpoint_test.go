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
	"testing"
	"time"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/camel"
	"github.com/stretchr/testify/assert"
)

// ==========================
//
// Duration converter
//
// ==========================

func TestNewEndpoint(t *testing.T) {
	ctx := camel.NewContextWithParent(nil)
	ctx.AddTypeConverter(camel.ToIntConverter)
	ctx.AddTypeConverter(camel.ToDurationConverter)
	ctx.Registry().Bind("http", NewComponent())

	endpoint, err := api.NewEndpointFromURI(ctx, "http://www.google.com/search?requestTimeout=25s")

	assert.NoError(t, err)
	assert.NotNil(t, endpoint)

	httpe, ok := endpoint.(*httpEndpoint)

	assert.NotEqual(t, ok, false)
	assert.NotNil(t, httpe)
	assert.Equal(t, "http", httpe.scheme)
	assert.Equal(t, "www.google.com", httpe.host)
	assert.Equal(t, 80, httpe.port)
	assert.Equal(t, "/search", httpe.path)
	assert.Equal(t, 10*time.Second, httpe.connectionTimeout)
	assert.Equal(t, 25*time.Second, httpe.requestTimeout)
	assert.Equal(t, ghttp.MethodGet, httpe.method)
	assert.Nil(t, httpe.transport)
	assert.Nil(t, httpe.client)
}
