package endpoint

import (
	"net/url"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/pkg/errors"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

const TAG = "endpoint"
const RemainingKey = "remaining"
const URIKey = "uri"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &Endpoint{}
	}
}

type Endpoint struct {
	api.Identifiable

	Identity   string                 `yaml:"id"`
	URI        string                 `yaml:"uri"`
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
}

func (e *Endpoint) ID() string {
	return e.Identity
}

func (e *Endpoint) Consumer(ctx api.Context, pid *actor.PID) (api.Consumer, error) {

	ep, err := e.create(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failure creating endpoint")
	}

	factory, ok := ep.(api.ConsumerFactory)
	if !ok {
		return nil, camelerrors.NotImplementedf("scheme %s does not implement consumer", ep.Component().Scheme())
	}

	consumer, err := factory.Consumer(pid)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating consumer")
	}

	return consumer, nil
}

func (e *Endpoint) Producer(ctx api.Context) (api.Producer, error) {

	ep, err := e.create(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failure creating endpoint")
	}

	factory, ok := ep.(api.ProducerFactory)
	if !ok {
		return nil, camelerrors.NotImplementedf("scheme %s does not implement producer", ep.Component().Scheme())
	}

	producer, err := factory.Producer()
	if err != nil {
		return nil, errors.Wrapf(err, "error creating producer")
	}

	return producer, nil
}

func (e *Endpoint) create(ctx api.Context) (api.Endpoint, error) {
	params := make(map[string]interface{})

	u, err := url.Parse(e.URI)
	if err != nil {
		return nil, err
	}

	for k, v := range u.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	for k, v := range e.Parameters {
		params[k] = v
	}

	r, _ := ctx.Properties().Expand(u.Opaque)

	params[RemainingKey] = r
	params[URIKey] = e.URI

	for k, v := range params {
		switch val := v.(type) {
		case string:
			v, _ = ctx.Properties().Expand(val)
			params[k] = v
		case []byte:
			v, _ = ctx.Properties().Expand(string(val))
			params[k] = v
		default:
			params[k] = val
		}
	}

	scheme, _ := ctx.Properties().Expand(u.Scheme)

	f, ok := components.Factories[scheme]
	if !ok {
		return nil, camelerrors.NotFoundf("not component for scheme %s", scheme)
	}

	c, err := f(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return c.Endpoint(params)
}
