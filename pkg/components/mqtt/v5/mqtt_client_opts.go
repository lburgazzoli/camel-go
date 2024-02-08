package v5

import (
	"github.com/eclipse/paho.golang/paho"
)

type OptionFn func(*paho.ClientConfig)

func WithSingleHandlerRouter(handler func(*paho.Publish)) OptionFn {
	return func(config *paho.ClientConfig) {
		config.Router = paho.NewStandardRouterWithDefault(handler)
	}
}
