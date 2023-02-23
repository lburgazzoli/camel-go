//go:build component_dapr || components_all

package dapr

import (
	"github.com/dapr/go-sdk/client"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"

	"github.com/mitchellh/mapstructure"

	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
)

const Scheme = "dapr"

func NewComponent(config map[string]interface{}) (api.Component, error) {
	c, err := client.NewClient()
	if err != nil {
		return nil, err
	}

	component := Component{
		id:     uuid.New(),
		scheme: Scheme,
		client: c,
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
	client client.Client
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
