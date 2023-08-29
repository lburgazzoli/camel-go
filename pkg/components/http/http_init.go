// //go:build components_http || components_all

package http

import (
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
)

func init() {
	components.Factories[SchemeHTTP] = func(context camel.Context, m map[string]interface{}) (camel.Component, error) {
		return newComponent(context, SchemeHTTP, m)
	}

	components.Factories[SchemeHTTPS] = func(context camel.Context, m map[string]interface{}) (camel.Component, error) {
		return newComponent(context, SchemeHTTPS, m)
	}
}
