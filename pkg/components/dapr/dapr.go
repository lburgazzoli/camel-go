//go:build component_dapr || components_all

package dapr

import (
	"github.com/dapr/go-sdk/client"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"

	"github.com/mitchellh/mapstructure"
)

const Scheme = "dapr"

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
	daprClient, err := client.NewClient()
	if err != nil {
		return nil, err
	}

	e := Endpoint{
		id:     uuid.New(),
		client: daprClient,
		config: c.config,
	}

	return &e, nil
}

// Endpoint ---
type Endpoint struct {
	id     string
	client client.Client
	config Config
}

func (e *Endpoint) ID() string {
	return e.id
}

func (e *Endpoint) Start() error {
	return nil
}

func (e *Endpoint) Stop() error {
	return nil
}
