package wasm

import (
	"github.com/dapr/go-sdk/client"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/mitchellh/mapstructure"
)

const Scheme = "wasm"

func NewComponent(config map[string]interface{}) (*Component, error) {
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
