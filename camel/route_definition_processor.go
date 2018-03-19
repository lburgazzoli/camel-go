package camel

// ==========================
//
// ProcessorDefinition
//
//    WORK IN PROGRESS
//
// ==========================

// NewProcessorDefinition --
func NewProcessorDefinition(context *Context) *ProcessorDefinition {
	def := ProcessorDefinition{}
	def.parent = nil
	def.child = nil
	def.context = context
	def.definitions = make([]definitionFactory, 0)

	return &def
}

// NewProcessorDefinitionWithParent --
func NewProcessorDefinitionWithParent(parent *ProcessorDefinition) *ProcessorDefinition {
	def := ProcessorDefinition{}
	def.parent = parent
	def.child = nil
	def.context = parent.context
	def.definitions = make([]definitionFactory, 0)

	return &def
}

// ProcessorDefinition --
type ProcessorDefinition struct {
	context     *Context
	definitions []definitionFactory
	child       *ProcessorDefinition
	parent      *ProcessorDefinition
}

// AddFactory --
func (definition *ProcessorDefinition) AddFactory(factory definitionFactory) *ProcessorDefinition {
	definition.definitions = append(definition.definitions, factory)

	return definition
}

// End --
func (definition *ProcessorDefinition) End() *ProcessorDefinition {
	return definition.parent
}

// To --
func (definition *ProcessorDefinition) To(uri string) *ProcessorDefinition {
	var err error
	var producer Producer
	var endpoint Endpoint

	if endpoint, err = definition.context.CreateEndpointFromURI(uri); err != nil {
		return nil
	}

	if producer, err = endpoint.CreateProducer(); err != nil {
		return nil
	}

	return definition.AddFactory(func(parent *Pipe) (*Pipe, Service) {
		p := producer.Pipe()

		parent.Subscribe(func(e *Exchange) {
			p.Publish(e)
		})

		return p, producer
	})
}
