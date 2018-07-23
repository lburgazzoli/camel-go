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

package api

import (
	"fmt"
	"net/url"

	zlog "github.com/rs/zerolog/log"
)

// HasContext --
type HasContext interface {
	Context() Context
}

// ContextAware --
type ContextAware interface {
	HasContext

	SetContext(context Context)
}

// Context --
type Context interface {
	Service

	// Registry
	Registry() LoadingRegistry

	// Type conversion
	AddTypeConverter(converter TypeConverter)
	//TypeConverters() []TypeConverter
	TypeConverter() TypeConverter

	// Routes
	AddRoute(route *Route)
	//Routes() []*Route

	// Services
	AddService(service Service) bool
	//Service() []Service
}

// NewEndpointFromURI --
func NewEndpointFromURI(context Context, uri string) (Endpoint, error) {
	var err error
	var endpointURL *url.URL
	var component Component
	var endpoint Endpoint

	if endpointURL, err = url.Parse(uri); err != nil {
		return nil, err
	}

	scheme := endpointURL.Scheme
	opts := make(map[string]interface{})
	vals := make(url.Values)

	if vals, err = url.ParseQuery(endpointURL.RawQuery); err != nil {
		return nil, err
	}

	for k, v := range vals {
		opts[k] = v[0]
	}

	if component, err = LookupComponent(context, scheme); err == nil {
		remaining := ""
		if endpointURL.Opaque != "" {
			if remaining, err = url.PathUnescape(endpointURL.Opaque); err != nil {
				return nil, err
			}
		} else {
			remaining = endpointURL.Host

			if endpointURL.Path != "" {
				path, err := url.PathUnescape(endpointURL.Path)
				if err != nil {
					return nil, err
				}

				remaining += path
			}
		}

		endpoint, err = component.CreateEndpoint(remaining, opts)
	}

	if err != nil {
		endpoint = nil
	}

	return endpoint, err
}

// LookupComponent --
func LookupComponent(context Context, scheme string) (Component, error) {
	var component Component
	var err error

	zlog.Info().Msgf("lookup component scheme: %s", scheme)

	// Every component should be  context registry
	if c, ok := context.Registry().Lookup(scheme); ok {
		zlog.Info().Msgf("scheme: %s, component: %v, error: %v", scheme, c, err)

		component, _ = c.(Component)
	}

	if component != nil {
		if ca, ok := component.(ContextAware); ok {
			ca.SetContext(context)
		}
	}

	if component == nil {
		err = fmt.Errorf("unable to find component whit scheme: %s", scheme)
	}

	return component, err
}
