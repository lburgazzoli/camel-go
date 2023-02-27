package from

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/endpoint"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/pkg/errors"
)

const TAG = "from"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &From{
			Endpoint: endpoint.Endpoint{
				Identity: uuid.New(),
			},
		}
	}
}

type From struct {
	endpoint.Endpoint `yaml:",inline"`

	Steps []processors.Step `yaml:"steps,omitempty"`
}

func (f *From) Reify(ctx api.Context) (string, error) {

	var last string

	for i := len(f.Steps) - 1; i >= 0; i-- {
		if last != "" {
			f.Steps[i].Next(last)
		}

		pid, err := f.Steps[i].Reify(ctx)
		if err != nil {
			return "", errors.Wrapf(err, "error creating step")
		}

		last = pid
	}

	if last != "" {
		f.Endpoint.Next(last)
	}

	consumer, err := f.Endpoint.Consumer(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "error creating consumer")
	}

	return consumer.ID(), ctx.Spawn(consumer)
}
