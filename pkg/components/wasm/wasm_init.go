package wasm

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
)

func init() {
	components.Factories[Scheme] = func(config map[string]interface{}) (api.Component, error) {
		return NewComponent(config)
	}
}
