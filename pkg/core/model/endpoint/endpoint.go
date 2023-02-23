package endpoint

import (
	"net/url"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/core/model"
)

const TAG = "endpoint"

func init() {
	model.Types[TAG] = func() interface{} {
		return &Endpoint{}
	}
}

type Endpoint struct {
	URL        url.URL                `yaml:"url,omitempty"`
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
}

func (e *Endpoint) Reify(_ api.Context) (actor.Actor, error) {

	params := make(map[string]interface{})

	for k, v := range e.URL.Query() {
		params[k] = v
	}
	for k, v := range e.Parameters {
		params[k] = v
	}

	return nil, errors.NotImplemented("")
}
