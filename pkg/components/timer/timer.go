////go:build components_timer || components_all

package timer

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
)

const (
	Scheme           = "timer"
	PropertiesPrefix = "camel.component." + Scheme

	AttributeTimerName       = "camel.apache.org/timer.name"
	AttributeTimerStarted    = "camel.apache.org/timer.started"
	AttributeTimerFiredCount = "camel.apache.org/timer.fired.count"
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
