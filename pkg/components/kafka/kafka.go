////go:build components_kafka || components_all

package kafka

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
)

const (
	Scheme           = "kafka"
	PropertiesPrefix = "camel.component." + Scheme

	AttributePartition = "camel.apache.org/kafka.partition"
	AttributeOffset    = "camel.apache.org/kafka.offset"
	AttributeTopic     = "camel.apache.org/kafka.topic"
	AttributeKey       = "camel.apache.org/kafka.key"
)

func NewComponent(ctx api.Context, _ map[string]interface{}) (api.Component, error) {
	component := Component{
		DefaultComponent: components.NewDefaultComponent(ctx, Scheme),
	}

	return &component, nil
}

type Component struct {
	components.DefaultComponent
}

func (c *Component) Endpoint(config api.Parameters) (api.Endpoint, error) {
	e := Endpoint{
		DefaultEndpoint: components.NewDefaultEndpoint(c),
	}

	props := c.Context().Properties().View(PropertiesPrefix).Parameters()
	for k, v := range config {
		props[k] = v
	}

	if _, err := c.Context().TypeConverter().Convert(&props, &e.config); err != nil {
		return nil, err
	}

	return &e, nil
}
