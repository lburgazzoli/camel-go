////go:build components_timer || components_all

package timer

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
)

const (
	Scheme = "timer"

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

	if _, err := c.Context().TypeConverter().Convert(&config, &e.config); err != nil {
		return nil, err
	}

	return &e, nil
}
