////go:build components_kafka || components_all

package kafka

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
	"github.com/mitchellh/mapstructure"
)

const Scheme = "kafka"

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
	e := Endpoint{
		DefaultEndpoint: components.NewDefaultEndpoint(c),
		config:          c.config,
	}

	return &e, nil
}

type Endpoint struct {
	components.DefaultEndpoint

	config Config
}

func (e *Endpoint) Start() error {
	return nil
}

func (e *Endpoint) Stop() error {
	return nil
}
