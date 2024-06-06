// //go:build components_mqtt_v3 || components_all

package v3

import (
	"time"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
)

const (
	Scheme           = "mqtt-v3"
	PropertiesPrefix = "camel.component." + Scheme

	AttributeMqttMessageID        = "camel.apache.org/mqtt.message.id"
	AttributeMqttMessageRetained  = "camel.apache.org/mqtt.message.retained"
	AttributeMqttMessageDuplicate = "camel.apache.org/mqtt.message.duplicate"
	AttributeMqttMessageQUOS      = "camel.apache.org/mqtt.message.qus"

	DefaultDisconnectTimeout = 250
	DefaultKeepAlive         = 2 * time.Second
	DefaultPingTimeout       = 1 * time.Second
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

	view, err := c.Context().Properties().View(PropertiesPrefix).Merge(config)
	if err != nil {
		return nil, err
	}

	params := view.Parameters()
	params = c.Context().Properties().ExpandAll(params)

	if _, err := c.Context().TypeConverter().Convert(&params, &e.config); err != nil {
		return nil, err
	}

	return &e, nil
}
