package from

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/endpoint"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/pkg/errors"
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

	current := api.OutputAware(&f.Endpoint)
	for i := range f.Steps {
		pid, err := f.Steps[i].Reify(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating step")
		}

		current.Next(pid)
		current = &f.Steps[i]
	}

	consumer, err := f.Endpoint.Consumer(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating consumer")
	}

	return ctx.SpawnFn(uuid.New(), func(c actor.Context) {
		switch c.Message().(type) {
		case *actor.Started:
			_ = consumer.Start()
		case *actor.Stopping:
			_ = consumer.Stop()
		}
	})
}
