package api

import (
	"reflect"
)

// ==========================
//
//
//
// ==========================

// Registry --
type Registry interface {
	Bind(name string, value interface{})

	Lookup(name string) (interface{}, bool)
	LookupAs(name string, expectedType reflect.Type) (interface{}, bool)

	Range(func(key string, value interface{}) bool)
	//LookupByType(expectedType reflect.Type) ([]interface{}, error)
}

// RegistryLoader --
type RegistryLoader interface {
	Load(name string) (interface{}, error)
}

// LoadingRegistry --
type LoadingRegistry interface {
	Registry

	AddLoader(loader RegistryLoader)
}
