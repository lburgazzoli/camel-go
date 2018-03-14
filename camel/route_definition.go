package camel

//
// NODE: this is an example
//

// ProcessorFactory --
type ProcessorFactory func(parent Processor) Processor

// RouteDefinition --
type RouteDefinition struct {
	From string
	To   string
}

// FromX --
func (definition *RouteDefinition) FromX(uri string) *ProcessorDefinition {
	return nil
}

// ToRoute --
func (definition *RouteDefinition) ToRoute(context *Context) (*Route, error) {
	fromEndpoint, _ := context.CreateEndpointFromURI(definition.From)
	toEndpoint, _ := context.CreateEndpointFromURI(definition.To)

	var err error
	var producer Producer
	var consumer Consumer

	t := NewPipeIn()
	f := NewPipeWithNext(t)

	if producer, err = toEndpoint.CreateProducer(t); err != nil {
		return nil, err
	}

	if consumer, err = fromEndpoint.CreateConsumer(f); err != nil {
		return nil, err
	}

	route := Route{}
	route.AddService(consumer)
	route.AddService(producer)

	return &route, nil
}

// ProcessorDefinition --
type ProcessorDefinition struct {
}

// ToX --
func (definition *ProcessorDefinition) ToX(uri string) *ProcessorDefinition {
	return nil
}
