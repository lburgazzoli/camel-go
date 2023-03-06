// //go:build components_log || components_all
package log

import (
	"context"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

type Endpoint struct {
	config Config
	components.DefaultEndpoint
}

func (e *Endpoint) Start(context.Context) error {
	return nil
}

func (e *Endpoint) Stop(context.Context) error {
	return nil
}

func (e *Endpoint) Producer() (api.Producer, error) {
	c := Producer{
		id:       uuid.New(),
		endpoint: e,
		logger:   e.Logger().Named(e.config.Remaining),
	}

	return &c, nil
}
