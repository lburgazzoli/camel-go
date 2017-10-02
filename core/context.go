package core

import (
	"github.com/lburgazzoli/camel-go/camel"
)

type defaultCamelContext struct {
	camel.Service

	name            string
	registryLoaders []camel.RegistryLoader
	components      map[string]camel.Component
}

// ==========================
//
// Initialize a camel context
//
// ==========================

// NewCamelContext --
func NewCamelContext() camel.Context {
	return &defaultCamelContext{
		Service:         NewService(),
		name:            "camel",
		registryLoaders: make([]camel.RegistryLoader, 0),
		components:      make(map[string]camel.Component),
	}
}

// NewCamelContextWithName --
func NewCamelContextWithName(name string) camel.Context {
	return &defaultCamelContext{
		Service:         NewService(),
		name:            name,
		registryLoaders: make([]camel.RegistryLoader, 0),
		components:      make(map[string]camel.Component),
	}
}

// ==========================
//
//
//
// ==========================

func (context *defaultCamelContext) AddRegistryLoader(loader camel.RegistryLoader) {
	context.registryLoaders = append(context.registryLoaders, loader)
}

func (context *defaultCamelContext) AddComponent(name string, component camel.Component) {
	context.components[name] = component
	context.components[name].SetContext(context)
}

func (context *defaultCamelContext) Component(name string) (camel.Component, error) {
	component, found := context.components[name]

	// check if the component has already been created or added to the context
	// component list
	if !found {
		for _, loader := range context.registryLoaders {
			component, err := loader.Load(name)

			if err != nil {
				return nil, err
			}

			if component == nil {
				continue
			}

			if _, ok := component.(camel.Component); !ok {
				// not a component
				continue
			}

			if component != nil {
				break
			}
		}

		if component != nil {
			context.AddComponent(name, component)
		}
	}

	return component, nil
}

// ==========================
//
// Lyfecycle
//
// ==========================

// Start --
func (context *defaultCamelContext) Start() {
}

// Stop --
func (context *defaultCamelContext) Stop() {
}
