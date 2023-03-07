package endpoint

import (
	"net/url"

	"github.com/pkg/errors"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
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
	api.WithOutputs

	Identity   string                 `yaml:"id"`
	URI        string                 `yaml:"uri"`
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
}

func (e *Endpoint) ID() string {
	return e.Identity
}

func (e *Endpoint) Consumer(ctx api.Context) (api.Consumer, error) {

	ep, err := e.create(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failure creating endpoint")
	}

	factory, ok := ep.(api.ConsumerFactory)
	if !ok {
		return nil, camelerrors.NotImplementedf("scheme %s does not implement consumer", ep.Component().Scheme())
	}

	consumer, err := factory.Consumer()
	if err != nil {
		return nil, errors.Wrapf(err, "error creating consumer")
	}

	for _, o := range e.Outputs() {
		next := o

		consumer.Next(next)
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

	for _, o := range e.Outputs() {
		next := o

		producer.Next(next)
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
			params[k] = ctx.Properties().String(v[0])
		}
	}

	for k, v := range e.Parameters {
		switch val := v.(type) {
		case string:
			params[k] = ctx.Properties().String(val)
		case []byte:
			params[k] = ctx.Properties().String(string(val))
		default:
			params[k] = val
		}
	}

	params["remaining"] = u.Opaque

	f, ok := components.Factories[u.Scheme]
	if !ok {
		return nil, camelerrors.NotFoundf("not component for scheme %s", u.Scheme)
	}

	c, err := f(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return c.Endpoint(params)
}
