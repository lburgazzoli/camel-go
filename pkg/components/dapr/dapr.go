////go:build components_dapr || components_all

package dapr

import (
	"context"

	"github.com/lburgazzoli/camel-go/pkg/util/serdes"

	"github.com/dapr/go-sdk/client"
	"github.com/lburgazzoli/camel-go/pkg/components"

	"github.com/lburgazzoli/camel-go/pkg/api"
)

const Scheme = "dapr"

func NewComponent(ctx api.Context, config map[string]interface{}) (api.Component, error) {

	component := Component{
		DefaultComponent: components.NewDefaultComponent(ctx, Scheme),
	}

	if err := serdes.DecodeStruct(&config, &component.config); err != nil {
		return nil, err
	}

	return &component, nil
}

type Component struct {
	config Config
	components.DefaultComponent
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
	client client.Client
	config Config
	components.DefaultEndpoint
}

func (e *Endpoint) Start(context.Context) error {
	return nil
}

func (e *Endpoint) Stop(context.Context) error {
	return nil
}
