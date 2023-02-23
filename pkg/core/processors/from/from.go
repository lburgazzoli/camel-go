package from

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/endpoint"
)

const TAG = "from"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &From{}
	}
}

type From struct {
	endpoint.Endpoint `yaml:",inline"`

	Steps []processors.Step `yaml:"steps,omitempty"`
}

func (f *From) Reify(ctx api.Context) (*actor.PID, error) {

	for i := range f.Steps {
		r, ok := f.Steps[i].T.(processors.Reifyable)
		if !ok {
			panic("internal error")
		}

		pid, err := r.Reify(ctx)
		if err != nil {
			return nil, err
		}

		f.Endpoint.Next(pid)
	}

	return f.Endpoint.Reify(ctx)
}
