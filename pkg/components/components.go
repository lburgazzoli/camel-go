package components

import "github.com/lburgazzoli/camel-go/pkg/api"

type ComponentFactory func(map[string]interface{}) (api.Component, error)

var Factories = make(map[string]ComponentFactory)
