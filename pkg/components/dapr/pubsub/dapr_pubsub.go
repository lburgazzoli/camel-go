// //go:build components_dapr_pubsub || components_all

package pubsub

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
	"github.com/lburgazzoli/camel-go/pkg/components/dapr"
)

const (
	Scheme = "dapr-pubsub"

	AttributeEventID     = "camel.apache.org/dapr.event.id"
	AttributePubSubName  = "camel.apache.org/dapr.pubsub.name"
	AttributePubSubTopic = "camel.apache.org/dapr.pubsub.topic"
)

func NewComponent(ctx api.Context, _ map[string]interface{}) (api.Component, error) {
	component := Component{}
	component.DefaultComponent = components.NewDefaultComponent(ctx, Scheme)
	component.s = NewService(dapr.Address(), component.Logger())

	return &component, nil
}

type Component struct {
	components.DefaultComponent

	s *Service
}

func (c *Component) Endpoint(config api.Parameters) (api.Endpoint, error) {
	e := Endpoint{
		DefaultEndpoint: components.NewDefaultEndpoint(c),
		s:               c.s,
	}

	if _, err := c.Context().TypeConverter().Convert(&config, &e.config); err != nil {
		return nil, err
	}

	return &e, nil
}
