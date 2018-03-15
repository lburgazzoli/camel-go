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

func (definition *ProcessorDefinition) addFactory(factory definitionFactory) *ProcessorDefinition {
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

	return definition.addFactory(func(parent *Pipe) (*Pipe, Service) {
		p := producer.Pipe()
		parent.Next(p)

		return p, producer
	})
}

// Process --
func (definition *ProcessorDefinition) Process(processor Processor) *ProcessorDefinition {
	return definition.addFactory(func(parent *Pipe) (*Pipe, Service) {
		return parent.Process(processor), nil
	})
}

// Filter --
func (definition *ProcessorDefinition) Filter() *FilterDefinition {
	filter := FilterDefinition{}
	filter.context = definition.context

	definition.child = &filter.ProcessorDefinition

	return &filter
}

// FilterWithPredicate --
func (definition *ProcessorDefinition) FilterWithPredicate(predicate Predicate) *ProcessorDefinition {
	return definition.Filter().WithPredicate(predicate)
}
