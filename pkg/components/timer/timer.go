////go:build components_timer || components_all

package timer

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
	"github.com/lburgazzoli/camel-go/pkg/util/serdes"
)

const (
	Scheme = "timer"

	AnnotationTimerName       = "camel.apache.org/timer.name"
	AnnotationTimerStarted    = "camel.apache.org/timer.started"
	AnnotationTimerFiredCount = "camel.apache.org/timer.fired.count"
)

func NewComponent(ctx api.Context, config map[string]interface{}) (api.Component, error) {
	component := Component{
		DefaultComponent: components.NewDefaultComponent(ctx, Scheme),
	}

	if err := serdes.Convert(&config, &component.config); err != nil {
		return nil, err
	}

	return &component, nil
}

type Component struct {
	components.DefaultComponent

	config Config
}

func (c *Component) Endpoint(params api.Parameters) (api.Endpoint, error) {
	e := Endpoint{
		DefaultEndpoint: components.NewDefaultEndpoint(c),
		config:          c.config,
	}

	return &e, nil
}
