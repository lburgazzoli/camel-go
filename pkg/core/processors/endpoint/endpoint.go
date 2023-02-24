package endpoint

import (
	"net/url"

	"github.com/lburgazzoli/camel-go/pkg/components"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

const TAG = "endpoint"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &Endpoint{}
	}
}

type Endpoint struct {
	api.Identifiable

	outputs *actor.PIDSet

	Identity   string                 `yaml:"id"`
	URL        url.URL                `yaml:"url,omitempty"`
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
}

func (e *Endpoint) ID() string {
	return e.Identity
}

func (e *Endpoint) Next(pid *actor.PID) {
	if e.outputs == nil {
		e.outputs = actor.NewPIDSet()
	}

	e.outputs.Add(pid)
}

func (e *Endpoint) Reify(_ api.Context) (*actor.PID, error) {

	params := make(map[string]interface{})

	for k, v := range e.URL.Query() {
		params[k] = v
	}
	for k, v := range e.Parameters {
		params[k] = v
	}

	id := e.Identity
	if id == "" {
		id = uuid.New()
	}

	_, ok := components.Factories[e.URL.Scheme]
	if !ok {
		return nil, errors.NotImplemented("")
	}

	return nil, errors.NotImplemented("")
}
