package core

import (
	"fmt"
	"github.com/lburgazzoli/camel-go/camel"
	"path"
	"plugin"
)

const (
	ComponentsDir       = "components"
	ComponentSymbolName = "CreateComponent"
)

type DefaultCamelContext struct {
	name string
}

func NewCamelContext() camel.Context {
	return &DefaultCamelContext{
		name: "camel",
	}
}

func NewCamelContextWithName(name string) camel.Context {
	return &DefaultCamelContext{
		name: name,
	}
}

func (context *DefaultCamelContext) Start() {
}

func (context *DefaultCamelContext) Stop() {
}

func (context *DefaultCamelContext) GetComponent(name string) (camel.Component, error) {
	pluginPath := path.Join(ComponentsDir, fmt.Sprintf("%s.so", name))

	plug, err := plugin.Open(pluginPath)
	if err != nil {
		fmt.Printf("failed to open plugin %s: %v\n", name, err)
		return nil, err
	}

	symbol, err := plug.Lookup(ComponentSymbolName)
	if err != nil {
		fmt.Printf("plugin %s does not export symbol \"%s\"\n", name, ComponentSymbolName)
		return nil, err
	}

	component := symbol.(func(camel.Context) camel.Component)(context)

	return component, nil
}
