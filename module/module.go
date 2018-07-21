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

package module

import (
	"fmt"
	"os"
	"plugin"

	zlog "github.com/rs/zerolog/log"
)

// LoadSymbol --
func LoadSymbol(path string, symbol string) (interface{}, error) {
	if path == "" {
		return nil, fmt.Errorf("Path is empty")
	}

	if symbol == "" {
		return nil, fmt.Errorf("Symbol is empty")
	}

	zlog.Debug().Msgf("try loading symbol \"%s\" from plugin %s", symbol, path)

	location := os.ExpandEnv(path)
	_, err := os.Stat(location)

	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	plug, err := plugin.Open(location)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin %s: %v", location, err)
	}

	answer, err := plug.Lookup(symbol)
	if err != nil {
		return nil, fmt.Errorf("plugin %s does not export symbol \"%s\"", location, symbol)
	}

	return answer, nil
}
