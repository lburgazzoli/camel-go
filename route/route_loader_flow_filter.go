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
	"fmt"
	"strings"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/module"
)

// FilterStepHandler --
func FilterStepHandler(step Step, route *RouteDefinition) (*RouteDefinition, error) {
	impl := struct {
		TypedStep

		Function string `yaml:"function"`
		Location string `yaml:"location"`
		Language string `yaml:"language"`
	}{}

	err := decodeStep("filter", step, &impl)
	if err != nil {
		return nil, err
	}

	if impl.Function == "" {
		return nil, fmt.Errorf("missing function: %s", impl.Function)
	}

	// if the language is not set, we assume it is "ref"
	if impl.Language == "" || impl.Language == "ref" {

		// check if the function is defined in an external
		// plugin
		if strings.HasPrefix(impl.Location, "file:") {
			location := strings.TrimPrefix(impl.Location, "file:")
			symbol, err := module.LoadSymbol(location, impl.Function)

			if err != nil {
				return nil, err
			}

			return route.Filter().Fn(symbol.(func(api.Exchange) bool)), nil
		}

		return route.Filter().Ref(impl.Function), nil
	}

	if impl.Language == "jsonpath" {
		return route.Filter().JSONPath(impl.Function), nil
	}

	return nil, fmt.Errorf("unsupported language: %s", impl.Language)
}
