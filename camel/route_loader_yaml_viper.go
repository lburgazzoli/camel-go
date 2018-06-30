package camel

import (
	"github.com/spf13/viper"
)

// ==========================
//
// ViperYAMLLoader
//
// ==========================

// ViperYAMLLoader --
type ViperYAMLLoader struct {
	loader FlowLoader
	viper  *viper.Viper
}

// ==========================
//
// Initialization
//
// ==========================

// NewViperYAMLLoader --
func NewViperYAMLLoader(viper *viper.Viper) RouteLoader {
	loader := ViperYAMLLoader{
		loader: NewFlowLoader(),
		viper:  viper,
	}

	return &loader
}

// ==========================
//
// Implementation
//
// ==========================

// Load --
func (loader *ViperYAMLLoader) Load() ([]Definition, error) {
	flows := make([]Flow, 0)
	err := loader.viper.UnmarshalKey("flows", &flows)

	if err != nil {
		return nil, err
	}

	return loader.loader.definition(flows)
}
