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

package route

import (
	"github.com/lburgazzoli/camel-go/api"
	"github.com/robertkrimen/otto"
)

// DefinitionWrapper --
type definitionWrapper struct {
	definitions []*RouteDefinition
}

// From --
func (target *definitionWrapper) from(uri string) *RouteDefinition {
	from := From(uri)
	target.definitions = append(target.definitions, from)

	return from
}

// LoadFromJS --
func LoadFromJS(context api.Context, js string) ([]*api.Route, error) {
	vm := otto.New()
	dw := new(definitionWrapper)
	dw.definitions = make([]*RouteDefinition, 0)

	if err := vm.Set("From", dw.from); err != nil {
		return nil, err
	}

	if _, err := vm.Run(js); err != nil {
		return nil, err
	}

	routes := make([]*api.Route, 0)

	for _, d := range dw.definitions {
		r, e := ToRoute(context, d)
		if e != nil {
			return nil, e
		}

		routes = append(routes, r)
	}

	return routes, nil
}
