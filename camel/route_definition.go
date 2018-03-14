package camel

//
// NODE: this is an example
//

// RouteDefinition --
type RouteDefinition struct {
	From string
	To   string
}

// ToRoute --
func (definition *RouteDefinition) ToRoute(context *Context) (*Route, error) {
	fromEndpoint, _ := context.CreateEndpointFromURI(definition.From)
	toEndpoint, _ := context.CreateEndpointFromURI(definition.To)

	var err error
	var producer Producer
	var consumer Consumer

	if producer, err = toEndpoint.CreateProducer(); err != nil {
		return nil, err
	}

	if consumer, err = fromEndpoint.CreateConsumer(); err != nil {
		return nil, err
	}

	producer.Pipe().In = make(chan *Exchange)
	consumer.Pipe().Next = producer.Pipe()

	route := Route{}
	route.AddService(consumer)
	route.AddService(producer)

	return &route, nil
}

// RouteDefinitionNg --
type RouteDefinitionNg struct {
	context *Context
}

// From --
func (definition *RouteDefinitionNg) From(uri string) *ProcessorDefinitionNg {
	return nil
}

// ProcessorDefinitionNg --
type ProcessorDefinitionNg struct {
	context *Context
}

// To --
func (definition *ProcessorDefinitionNg) To(uri string) *ProcessorDefinitionNg {
	var err error
	var producer Producer
	var endpoint Endpoint

	if endpoint, err = definition.context.CreateEndpointFromURI(uri); err != nil {
		return nil
	}

	if producer, err = endpoint.CreateProducer(); err != nil {
		return nil
	}

	producer.Pipe()

	return nil
}
