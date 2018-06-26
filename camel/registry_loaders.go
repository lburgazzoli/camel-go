package camel

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

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
func NewPluginRegistryLoader(searchPath string) RegistryLoader {
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
		var symbol interface{}
		var err error
		var pluginPath string

		// first scan all the plugins to find one that export the
		// name as symbol
		symbol, err = loader.scanForSymbol(name)
		if err != nil {
			return nil, err
		}

		if symbol != nil {
			result = symbol
		}

		if result == nil {
			// then lookup a factory
			pluginPath = path.Join(loader.searchPath, fmt.Sprintf("%s.so", name))
			symbol, err = LoadSymbol(pluginPath, "Create")

			if err != nil {
				zlog.Warn().Msgf("plugin %s does not export symbol \"Create\"", name)
				return nil, err
			}

			// Load the object from
			result = symbol.(func() interface{})()
		}

		loader.cache[name] = result
	}

	return result, nil
}

// scanForSymbol --
func (loader *pluginRegistryLoader) scanForSymbol(name string) (interface{}, error) {
	files, err := ioutil.ReadDir(loader.searchPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		ext := path.Ext(file.Name())
		if ext == ".so" {
			symbol, err := LoadSymbol(path.Join(loader.searchPath, file.Name()), name)

			if err != nil {
				return nil, err
			}

			return symbol, nil
		}
	}

	return nil, nil
}
