package camel

import (
	"net/url"
)

//
// NODE: this is an example
//

// RouteDefinition --
type RouteDefinition struct {
	From string
	To   string
}

// CreateRoute --
func (definition *RouteDefinition) CreateRoute(context *Context) (*Route, error) {
	var err error
	var fromURL *url.URL
	var toURL *url.URL
	var fromComponent Component
	var toComponent Component
	var fromEndpoint Endpoint
	var toEndpoint Endpoint
	var producer Producer
	var consumer Consumer

	if fromURL, err = url.Parse(definition.From); err != nil {
		return nil, err
	}
	if toURL, err = url.Parse(definition.To); err != nil {
		return nil, err
	}
	if fromComponent, err = context.Component(fromURL.Scheme); err != nil {
		return nil, err
	}
	if toComponent, err = context.Component(toURL.Scheme); err != nil {
		return nil, err
	}

	fromVals, _ := url.ParseQuery(fromURL.RawQuery)
	fromOpts := make(map[string]interface{})
	for k, v := range fromVals {
		fromOpts[k] = v[0]
	}

	if fromEndpoint, err = fromComponent.CreateEndpoint("", fromOpts); err != nil {
		return nil, err
	}

	toVals, _ := url.ParseQuery(toURL.RawQuery)
	toOpts := make(map[string]interface{})
	for k, v := range toVals {
		toOpts[k] = v[0]
	}

	if toEndpoint, err = toComponent.CreateEndpoint("test", toOpts); err != nil {
		return nil, err
	}
	if producer, err = toEndpoint.CreateProducer(); err != nil {
		return nil, err
	}
	if consumer, err = fromEndpoint.CreateConsumer(producer); err != nil {
		return nil, err
	}

	route := Route{}
	route.AddService(consumer)
	route.AddService(producer)

	return &route, nil
}
