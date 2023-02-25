package from

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/endpoint"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/pkg/errors"
)

const TAG = "to"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &To{}
	}
}

type To struct {
	endpoint.Endpoint `yaml:",inline"`
}

func (t *To) Reify(ctx api.Context) (*actor.PID, error) {
	producer, err := t.Endpoint.Producer(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating consumer")
	}

	return ctx.Spawn(uuid.New(), producer)
}
