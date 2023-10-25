// //go:build components_dapr_pubsub || components_all

package pubsub

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
)

const (
	Scheme           = "dapr-pubsub"
	PropertiesPrefix = "camel.component." + Scheme

	AttributeEventID     = "camel.apache.org/dapr.event.id"
	AttributePubSubName  = "camel.apache.org/dapr.pubsub.name"
	AttributePubSubTopic = "camel.apache.org/dapr.pubsub.topic"
)

func NewComponent(ctx api.Context, _ map[string]interface{}) (api.Component, error) {
	component := Component{}
	component.DefaultComponent = components.NewDefaultComponent(ctx, Scheme)

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

	props = c.Context().Properties().ExpandAll(props)

	if _, err := c.Context().TypeConverter().Convert(&props, &e.config); err != nil {
		return nil, err
	}

	return &e, nil
}
