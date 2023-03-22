package from

import (
	"context"

	camel "github.com/lburgazzoli/camel-go/pkg/api"

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

func (f *From) Reify(ctx context.Context) (string, error) {
	camelContext := camel.GetContext(ctx)

	if err := processors.ReifySteps(ctx, &f.Endpoint, f.Steps); err != nil {
		return "", errors.Wrapf(err, "error creating from steps")
	}

	consumer, err := f.Endpoint.Consumer(camelContext)
	if err != nil {
		return "", errors.Wrapf(err, "error creating consumer")
	}

	return consumer.ID(), camelContext.Spawn(consumer)
}
