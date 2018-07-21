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
	"fmt"
	"os"
	"path"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/module"

	zlog "github.com/rs/zerolog/log"
)

// ==========================
//
// PluginRegistryLoader
//
//     Use Go's plugins to load objects
//
// ==========================

// NewPluginRegistryLoader --
func NewPluginRegistryLoader(searchPath string) api.RegistryLoader {
	return &pluginRegistryLoader{
		cache:      make(map[string]interface{}),
		searchPath: os.ExpandEnv(searchPath),
	}
}

type pluginRegistryLoader struct {
	cache      map[string]interface{}
	searchPath string
}

// Status --
func (loader *pluginRegistryLoader) Status() api.ServiceStatus {
	return api.ServiceStatusSTARTED
}

// Start --
func (loader *pluginRegistryLoader) Start() {
	// maybe here we should scan the search path to pre instantiate objects
}

// Stop --
func (loader *pluginRegistryLoader) Stop() {
}

// Load --
func (loader *pluginRegistryLoader) Load(name string) (interface{}, error) {
	var result, found = loader.cache[name]

	if !found {
		// then lookup a factory
		pluginPath := path.Join(loader.searchPath, fmt.Sprintf("%s.so", name))
		symbol, err := module.LoadSymbol(pluginPath, "Create")

		if err != nil {
			zlog.Warn().Msgf("plugin %s does not export symbol \"Create\"", name)
			return nil, err
		}

		if symbol == nil {
			return nil, nil
		}

		// Load the object from
		result = symbol.(func() interface{})()

		loader.cache[name] = result
	}

	return result, nil
}

/*
// scanForSymbol --
func (loader *pluginRegistryLoader) scanForSymbol(name string) (interface{}, error) {
	files, err := ioutil.ReadDir(loader.searchPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		ext := path.Ext(file.Name())
		if ext == ".so" {
			symbol, err := module.LoadSymbol(path.Join(loader.searchPath, file.Name()), name)

			if err != nil {
				return nil, err
			}

			return symbol, nil
		}
	}

	return nil, nil
}
*/
