package camel

import (
	"reflect"

	"github.com/lburgazzoli/camel-go/types"
)

// ==========================
//
// Initialize Registry
//
// ==========================

// NewRegistry --
func NewRegistry(converter types.TypeConverter) *Registry {
	return &Registry{
		converter: converter,
		local:     make(map[string]interface{}),
		loaders:   make([]RegistryLoader, 0),
	}
}

// ==========================
//
// Registry
//
// ==========================

// Registry --
type Registry struct {
	converter types.TypeConverter
	local     map[string]interface{}
	loaders   []RegistryLoader
}

// AddLoader --
func (registry *Registry) AddLoader(loader RegistryLoader) {
	registry.loaders = append(registry.loaders, loader)
}

// Bind --
func (registry *Registry) Bind(name string, value interface{}) {
	old, found := registry.local[name]
	if found {
		if service, ok := old.(Service); ok {
			service.Stop()
		}
	}

	registry.local[name] = value
}

// Lookup --
func (registry *Registry) Lookup(name string) (interface{}, error) {
	var value, found = registry.local[name]

	// check if the value has already been created
	if !found {
		for _, loader := range registry.loaders {
			value, err := loader.Load(name)

			if err != nil {
				return nil, err
			}

			if value == nil {
				continue
			}

			if value != nil {
				break
			}
		}
	}

	return value, nil
}

// LookupAs --
func (registry *Registry) LookupAs(name string, expectedType reflect.Type) (interface{}, error) {
	value, err := registry.Lookup(name)

	// check if the value has already been created
	if err != nil {
		return nil, err
	}

	// Convert to the expected type
	return registry.converter(value, expectedType)
}

// LookupAsOf --
func (registry *Registry) LookupAsOf(name string, expectedType interface{}) (interface{}, error) {
	return registry.LookupAs(name, reflect.TypeOf(expectedType))
}
