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
	"strings"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/module"
)

// ProcessStepHandler --
func ProcessStepHandler(step Step, route *RouteDefinition) (*RouteDefinition, error) {
	impl := struct {
		TypedStep

		Ref      string `yaml:"ref"`
		Location string `yaml:"location"`
	}{}

	err := decodeStep("process", step, &impl)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(impl.Location, "file:") {
		location := strings.TrimPrefix(impl.Location, "file:")
		symbol, err := module.LoadSymbol(location, impl.Ref)

		if err != nil {
			return nil, err
		}

		return route.Process().Fn(symbol.(func(api.Exchange))), nil
	}

	return route.Process().Ref(impl.Ref), nil
}
