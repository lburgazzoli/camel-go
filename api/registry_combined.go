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
)

// NewCombinedRegistry --
func NewCombinedRegistry(registry LoadingRegistry, registries ...Registry) LoadingRegistry {
	c := CombinedRegistry{
		root:    registry,
		parents: make([]Registry, 0),
	}

	for _, r := range registries {
		c.parents = append(c.parents, r)
	}

	return &c
}

// CombinedRegistry --
type CombinedRegistry struct {
	root    LoadingRegistry
	parents []Registry
}

// AddLoader --
func (registry *CombinedRegistry) AddLoader(loader RegistryLoader) {
	registry.root.AddLoader(loader)
}

// Bind --
func (registry *CombinedRegistry) Bind(name string, value interface{}) {
	registry.root.Bind(name, value)
}

// Lookup --
func (registry *CombinedRegistry) Lookup(name string) (interface{}, bool) {
	var answer interface{}
	var found bool

	answer, found = registry.root.Lookup(name)
	if !found {
		for _, r := range registry.parents {
			answer, found = r.Lookup(name)

			if found {
				break
			}
		}
	}

	return answer, found
}

// LookupAs --
func (registry *CombinedRegistry) LookupAs(name string, asType reflect.Type) (interface{}, bool) {
	var answer interface{}
	var found bool

	answer, found = registry.root.LookupAs(name, asType)
	if !found {
		for _, r := range registry.parents {
			answer, found = r.LookupAs(name, asType)

			if found {
				break
			}
		}
	}

	return answer, found
}

// Range --
func (registry *CombinedRegistry) Range(f func(key string, value interface{}) bool) {
	registry.root.Range(f)

	// TODO: need to be revisited
	for _, r := range registry.parents {
		r.Range(f)
	}
}

// ForEach --
func (registry *CombinedRegistry) ForEach(f func(key string, value interface{})) {
	registry.root.Range(func(key string, value interface{}) bool {
		f(key, value)

		return true
	})

	// TODO: need to be revisited
	for _, r := range registry.parents {
		r.Range(func(key string, value interface{}) bool {
			f(key, value)

			return true
		})
	}
}
