package core

import (
    "plugin"
	"github.com/lburgazzoli/camel-go/camel"
	"os"
	"path"
	"fmt"
)

const (
	ComponentsDir = "components"
	ComponentSymbolName = "Component"
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

	component, ok := symbol.(camel.Component)
	if !ok {
		fmt.Printf("Symbol from %s does not implement Component interface\n", name)
		return nil, nil
	}

	if err := component.Init(context); err != nil {
		fmt.Printf("%s initialization failed: %v\n", name, err)
		return nil, err
	}

	return component, nil
}
