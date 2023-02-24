//go:build components_timer || components_all

package timer

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"

	"github.com/mitchellh/mapstructure"
)

const Scheme = "timer"

func NewComponent(config map[string]interface{}) (api.Component, error) {
	component := Component{
		id:     uuid.New(),
		scheme: Scheme,
	}

	if err := mapstructure.Decode(config, &component.config); err != nil {
		return nil, err
	}

	return &component, nil
}

// Component ---
type Component struct {
	id     string
	scheme string
	config Config
}

func (c *Component) ID() string {
	return c.id
}

func (c *Component) Scheme() string {
	return c.scheme
}

func (c *Component) Endpoint(api.Parameters) (api.Endpoint, error) {
	e := Endpoint{
		id:     uuid.New(),
		config: c.config,
	}

	return &e, nil
}
