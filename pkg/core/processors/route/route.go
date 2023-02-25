package route

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/from"
)

const TAG = "route"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &Route{}
	}
}

type Route struct {
	api.Identifiable

	Identity string    `yaml:"id,omitempty"`
	Group    string    `yaml:"group,omitempty"`
	From     from.From `yaml:"from"`
}

func (r *Route) ID() string {
	return r.Identity
}

func (r *Route) Reify(ctx api.Context) (*actor.PID, error) {
	return r.From.Reify(ctx)
}
