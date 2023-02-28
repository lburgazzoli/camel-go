////go:build components_log || components_all

package log

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
	"github.com/lburgazzoli/camel-go/pkg/util/serdes"
)

const Scheme = "log"

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
	e := Endpoint{
		DefaultEndpoint: components.NewDefaultEndpoint(c),
		config:          c.config,
	}

	return &e, nil
}
