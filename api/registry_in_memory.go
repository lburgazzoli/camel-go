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
	"reflect"
	"sync"
)

// ==========================
//
//
//
// ==========================

// NewInMemoryRegistry --
func NewInMemoryRegistry(typeConverter TypeConverter) Registry {
	r := InMemoryRegistry{
		typeConverter: typeConverter,
	}

	return &r
}

// InMemoryRegistry --
type InMemoryRegistry struct {
	typeConverter TypeConverter
	data          sync.Map
}

// Bind --
func (registry *InMemoryRegistry) Bind(name string, value interface{}) {
	registry.data.Store(name, value)
}

// Lookup --
func (registry *InMemoryRegistry) Lookup(name string) (interface{}, bool) {
	return registry.data.Load(name)
}

// LookupAs --
func (registry *InMemoryRegistry) LookupAs(name string, asType reflect.Type) (interface{}, bool) {
	answer, found := registry.Lookup(name)

	if found {
		result, err := registry.typeConverter(answer, asType)

		if err != nil {
			return nil, false
		}

		return result, true
	}

	return answer, true
}

// Range --
func (registry *InMemoryRegistry) Range(f func(key string, value interface{}) bool) {
	registry.data.Range(func(key, value interface{}) bool {
		return f(key.(string), value)
	})
}

// ForEach --
func (registry *InMemoryRegistry) ForEach(f func(key string, value interface{})) {
	registry.data.Range(func(key, value interface{}) bool {
		f(key.(string), value)

		return true
	})
}
