////go:build components_dapr || components_all

package dapr

import (
	"github.com/dapr/go-sdk/client"
	"github.com/lburgazzoli/camel-go/pkg/components"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/mitchellh/mapstructure"
)

const Scheme = "dapr"

func NewComponent(ctx api.Context, config map[string]interface{}) (api.Component, error) {

	component := Component{
		DefaultComponent: components.NewDefaultComponent(ctx, Scheme),
	}

	if err := mapstructure.Decode(config, &component.config); err != nil {
		return nil, err
	}

	return &component, nil
}

type Component struct {
	components.DefaultComponent

	config Config
}

func (c *Component) Endpoint(api.Parameters) (api.Endpoint, error) {
	daprClient, err := client.NewClient()
	if err != nil {
		return nil, err
	}

	e := Endpoint{
		DefaultEndpoint: components.NewDefaultEndpoint(c),
		client:          daprClient,
		config:          c.config,
	}

	return &e, nil
}

type Endpoint struct {
	components.DefaultEndpoint

	client client.Client
	config Config
}

func (e *Endpoint) Start() error {
	return nil
}

func (e *Endpoint) Stop() error {
	return nil
}
