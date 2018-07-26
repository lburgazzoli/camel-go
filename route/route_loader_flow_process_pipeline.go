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

// ProcessPipelineStepHandler --
func ProcessPipelineStepHandler(step Step, route *RouteDefinition) (*RouteDefinition, error) {
	impl := struct {
		TypedStep

		Processors []struct {
			Ref      string `yaml:"ref"`
			Location string `yaml:"location"`
		}
	}{}

	err := decodeStep("pipeline", step, &impl)
	if err != nil {
		return nil, err
	}

	pipeline := route.Pipeline()

	for _, p := range impl.Processors {
		if strings.HasPrefix(p.Location, "file:") {
			location := strings.TrimPrefix(p.Location, "file:")
			name := p.Ref

			if name == "" {
				name = "Create"
			}

			symbol, err := module.LoadSymbol(location, p.Ref)

			if err != nil {
				return nil, err
			}

			pipeline.Fn(symbol.(func(api.Exchange)))
		} else if p.Ref != "" {
			pipeline.Ref(p.Ref)
		}
	}

	return pipeline.End(), nil
}
