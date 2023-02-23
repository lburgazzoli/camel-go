package route

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/endpoint"
)

const TAG = "route"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &Route{}
	}
}

type Route struct {
	ID    string            `yaml:"id,omitempty"`
	Group string            `yaml:"group,omitempty"`
	From  endpoint.Endpoint `yaml:"from"`
}

func (r *Route) Reify(_ api.Context) (*actor.PID, error) {
	return nil, errors.NotImplemented("")
}
