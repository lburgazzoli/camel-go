package api

import (
	"reflect"
	"sync"
)

// ==========================
//
//
//
// ==========================

// NewInMemoryRegistry --
func NewInMemoryRegistry(typeConverter TypeConverter) Registry {
	r := InMemoryRegistry{
		typeConverter: typeConverter,
	}

	return &r
}

// InMemoryRegistry --
type InMemoryRegistry struct {
	typeConverter TypeConverter
	data          sync.Map
}

// Bind --
func (registry *InMemoryRegistry) Bind(name string, value interface{}) {
	registry.data.Store(name, value)
}

// Lookup --
func (registry *InMemoryRegistry) Lookup(name string) (interface{}, bool) {
	return registry.data.Load(name)
}

// LookupAs --
func (registry *InMemoryRegistry) LookupAs(name string, asType reflect.Type) (interface{}, bool) {
	answer, found := registry.Lookup(name)

	if found {
		result, err := registry.typeConverter(answer, asType)

		if err != nil {
			return nil, false
		}

		return result, true
	}

	return answer, true
}

// Range --
func (registry *InMemoryRegistry) Range(f func(key string, value interface{}) bool) {
	registry.data.Range(func(key, value interface{}) bool {
		return f(key.(string), value)
	})
}
