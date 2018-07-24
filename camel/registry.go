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

package camel

import (
	"reflect"

	"github.com/lburgazzoli/camel-go/api"
)

// ==========================
//
// Initialize Registry
//
// ==========================

// NewRegistry --
func NewRegistry(converter api.TypeConverter) api.LoadingRegistry {
	return &defaultRegistry{
		converter: converter,
		local:     api.NewInMemoryRegistry(converter),
		loaders:   make([]api.RegistryLoader, 0),
	}
}

// ==========================
//
// defaultRegistry
//
// ==========================

// defaultRegistry --
type defaultRegistry struct {
	converter api.TypeConverter
	local     api.Registry
	loaders   []api.RegistryLoader
}

// AddLoader --
func (registry *defaultRegistry) AddLoader(loader api.RegistryLoader) {
	registry.loaders = append(registry.loaders, loader)
}

// Bind --
func (registry *defaultRegistry) Bind(name string, value interface{}) {
	old, found := registry.local.Lookup(name)
	if found {
		if service, ok := old.(api.Service); ok {
			service.Stop()
		}
	}

	registry.local.Bind(name, value)
}

// Lookup --
func (registry *defaultRegistry) Lookup(name string) (interface{}, bool) {
	var value interface{}
	var found bool
	var err error

	value, found = registry.local.Lookup(name)

	if !found {
		for _, loader := range registry.loaders {
			value, err = loader.Load(name)

			if err != nil {
				return nil, false
			}

			if value != nil {
				return value, true
			}
		}
	}

	return value, value != nil
}

// LookupAs --
func (registry *defaultRegistry) LookupAs(name string, expectedType reflect.Type) (interface{}, bool) {
	var value interface{}
	var found bool
	var err error

	if value, found = registry.Lookup(name); found {
		// Convert to the expected type
		value, err = registry.converter(value, expectedType)
		if err != nil {
			return nil, false
		}
	}

	return value, value != nil
}

// Range --
func (registry *defaultRegistry) Range(f func(key string, value interface{}) bool) {
	registry.local.Range(f)
}

// Range --
func (registry *defaultRegistry) ForEach(f func(key string, value interface{})) {
	registry.local.Range(func(key string, value interface{}) bool {
		f(key, value)

		return true
	})
}
