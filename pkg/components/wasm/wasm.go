//go:build component_wasm || components_all

package wasm

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"

	"github.com/mitchellh/mapstructure"

	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
)

const Scheme = "wasm"

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
	return nil, camelerrors.NotImplementedf("Endpoint for scheme %s not implemented", c.Scheme())
}
