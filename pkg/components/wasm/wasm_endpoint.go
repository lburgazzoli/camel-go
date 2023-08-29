// //go:build components_wasm || components_all
package wasm

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
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
	return nil, camelerrors.NotImplemented("TODO")
}

func (e *Endpoint) Consumer(_ *actor.PID) (api.Consumer, error) {
	return nil, camelerrors.NotImplementedf("TODO")
}
