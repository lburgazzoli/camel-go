package from

import (
	"context"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/endpoint"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/pkg/errors"
)

const TAG = "to"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &To{
			Endpoint: endpoint.Endpoint{
				Identity: uuid.New(),
			},
		}
	}
}

type To struct {
	endpoint.Endpoint `yaml:",inline"`
}

func (t *To) Reify(_ context.Context, camelContext camel.Context) (string, error) {
	producer, err := t.Endpoint.Producer(camelContext)
	if err != nil {
		return "", errors.Wrapf(err, "error creating consumer")
	}

	return producer.ID(), camelContext.Spawn(producer)
}