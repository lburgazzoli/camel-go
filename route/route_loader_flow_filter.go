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
	"github.com/lburgazzoli/camel-go/module"
	"github.com/mitchellh/mapstructure"
	zlog "github.com/rs/zerolog/log"
)

// FilterStep --
type FilterStep struct {
	TypedStep

	Ref      string `yaml:"ref"`
	Location string `yaml:"location"`
}

// FilterStepHandler --
func FilterStepHandler(step Step, route *RouteDefinition) (*RouteDefinition, error) {
	var impl FilterStep

	// not really needed, added for testing purpose
	err := mapstructure.Decode(step, &impl)
	if err != nil {
		return nil, err
	}

	zlog.Debug().Msgf("handle filter: step=<%v>, impl=<%+v>", step, impl)

	if impl.Location != "" {
		symbol, err := module.LoadSymbol(impl.Location, impl.Ref)
		if err != nil {
			return nil, err
		}

		return route.Filter().Fn(symbol.(func(api.Exchange) bool)), nil
	}

	return route.Filter().Ref(impl.Ref), nil
}
