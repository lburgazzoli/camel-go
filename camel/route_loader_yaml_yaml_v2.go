package camel

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	zlog "github.com/rs/zerolog/log"
)

// ==========================
//
// YAMLLoader
//
// ==========================

// YAMLLoader --
type YAMLLoader struct {
	loader FlowLoader
	data   []byte
}

// ==========================
//
// Initialization
//
// ==========================

// NewYAMLLoader --
func NewYAMLLoader(data []byte) RouteLoader {
	loader := YAMLLoader{
		loader: NewFlowLoader(),
		data:   data,
	}

	return &loader
}

// ==========================
//
// Implementation
//
// ==========================

// Load --
func (loader *YAMLLoader) Load() ([]Definition, error) {
	integration := Integration{}
	err := yaml.Unmarshal([]byte(loader.data), &integration)

	if err != nil {
		return nil, err
	}

	return loader.loader.definition(integration.Flows)
}

// ==========================
//
// Helpers
//
// ==========================

// LoadRouteFromYAMLFile --
func LoadRouteFromYAMLFile(path string) ([]Definition, error) {
	zlog.Debug().Msgf("Loading routes from:  %s", path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	loader := NewYAMLLoader(data)
	return loader.Load()
}
