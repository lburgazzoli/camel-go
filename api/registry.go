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

// ==========================
//
//
//
// ==========================

// Registry --
type Registry interface {
	Bind(name string, value interface{})

	Lookup(name string) (interface{}, bool)
	LookupAs(name string, expectedType reflect.Type) (interface{}, bool)

	// Range calls f sequentially for each key and value present in the registry.
	// If f returns false, range stops the iteration.
	Range(func(key string, value interface{}) bool)
	ForEach(func(key string, value interface{}))
	//LookupByType(expectedType reflect.Type) ([]interface{}, error)
}

// RegistryLoader --
type RegistryLoader interface {
	Load(name string) (interface{}, error)
}

// LoadingRegistry --
type LoadingRegistry interface {
	Registry

	AddLoader(loader RegistryLoader)
}
