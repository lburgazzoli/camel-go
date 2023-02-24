package route

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/endpoint"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

const TAG = "route"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &Route{}
	}
}

type Route struct {
	api.Identifiable

	Identity string            `yaml:"id,omitempty"`
	Group    string            `yaml:"group,omitempty"`
	From     endpoint.Endpoint `yaml:"from"`
}

func (r *Route) ID() string {
	return r.Identity
}

func (r *Route) Reify(_ api.Context) (*actor.PID, error) {
	id := r.Identity
	if id == "" {
		id = uuid.New()
	}

	return nil, errors.NotImplemented("")
}
