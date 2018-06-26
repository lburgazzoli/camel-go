package camel

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

// ==========================
//
// PluginRegistryLoader
//
//     Use Go's plugins to load objects
//
// ==========================

// NewPluginRegistryLoader --
func NewPluginRegistryLoader(searchPath string) RegistryLoader {
	return &pluginRegistryLoader{
		cache:      make(map[string]interface{}),
		searchPath: searchPath,
	}
}

type pluginRegistryLoader struct {
	cache      map[string]interface{}
	searchPath string
}

// Status --
func (loader *pluginRegistryLoader) Status() ServiceStatus {
	return ServiceStatusSTARTED
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
		pluginPath := path.Join(loader.searchPath, fmt.Sprintf("%s.so", name))
		symbol, err := LoadSymbol(pluginPath, "Create")

		if err != nil {
			log.Printf("plugin %s does not export symbol \"Create\"\n", name)
			return nil, err
		}

		// Load the object from
		result = symbol.(func() interface{})()

		loader.cache[name] = result
	}

	return result, nil
}

// LoadAll --
func (loader *pluginRegistryLoader) LoadAll() ([]interface{}, error) {

	files, err := ioutil.ReadDir(loader.searchPath)
	if err != nil {
		return nil, err
	}

	answer := make([]interface{}, 0)
	for _, file := range files {
		ext := path.Ext(file.Name())
		if ext == ".so" {
			name := strings.TrimSuffix(file.Name(), ext)

			value, err := loader.Load(name)
			if err != nil {
				return nil, err
			}

			answer = append(answer, value)
		}
	}

	return answer, nil
}
