package api

import (
	"reflect"
)

// NewCombinedRegistry --
func NewCombinedRegistry(registry LoadingRegistry, registries ...Registry) LoadingRegistry {
	c := CombinedRegistry{
		root:    registry,
		parents: make([]Registry, 0),
	}

	for _, r := range registries {
		c.parents = append(c.parents, r)
	}

	return &c
}

// CombinedRegistry --
type CombinedRegistry struct {
	root    LoadingRegistry
	parents []Registry
}

// AddLoader --
func (registry *CombinedRegistry) AddLoader(loader RegistryLoader) {
	registry.root.AddLoader(loader)
}

// Bind --
func (registry *CombinedRegistry) Bind(name string, value interface{}) {
	registry.root.Bind(name, value)
}

// Lookup --
func (registry *CombinedRegistry) Lookup(name string) (interface{}, bool) {
	var answer interface{}
	var found bool

	answer, found = registry.root.Lookup(name)
	if !found {
		for _, r := range registry.parents {
			answer, found = r.Lookup(name)

			if found {
				break
			}
		}
	}

	return answer, found
}

// LookupAs --
func (registry *CombinedRegistry) LookupAs(name string, asType reflect.Type) (interface{}, bool) {
	var answer interface{}
	var found bool

	answer, found = registry.root.LookupAs(name, asType)
	if !found {
		for _, r := range registry.parents {
			answer, found = r.LookupAs(name, asType)

			if found {
				break
			}
		}
	}

	return answer, found
}

// Range --
func (registry *CombinedRegistry) Range(f func(key string, value interface{}) bool) {
	registry.root.Range(f)

	// TODO: need to be revisited
	for _, r := range registry.parents {
		r.Range(f)
	}
}
