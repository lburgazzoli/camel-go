////go:build components_http || components_all

package http

import (
	"net/http"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
)

const (
	Scheme           = "http"
	PropertiesPrefix = "camel.component." + Scheme

	AttributeStatusMessage = "camel.apache.org/http.status.message"
	AttributeStatusCode    = "camel.apache.org/http.status.code"
	AttributeProto         = "camel.apache.org/http.proto"
	AttributeContentLength = "camel.apache.org/http.content-length"
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

	if e.config.Method == "" {
		e.config.Method = http.MethodGet
	}

	return &e, nil
}
