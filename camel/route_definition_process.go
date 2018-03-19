package camel

// ==========================
//
// ProcessDefinition
//
//    WORK IN PROGRESS
//
// ==========================

// Process --
func (definition *ProcessorDefinition) Process() *ProcessDefinition {
	process := ProcessDefinition{}
	process.context = definition.context
	process.parent = definition

	definition.child = &process.ProcessorDefinition

	return &process
}

// ==========================
//
//
//
//
//
// ==========================

// ProcessDefinition --
type ProcessDefinition struct {
	ProcessorDefinition
	processor Processor
}

// P --
func (definition *ProcessDefinition) P(processor Processor) *ProcessorDefinition {
	definition.processor = processor
	definition.child = NewProcessorDefinitionWithParent(&definition.ProcessorDefinition)

	definition.AddFactory(func(parent *Pipe) (*Pipe, Service) {
		return NewProcessorPipe(parent, processor), nil
	})

	return definition.child
}

// Fn --
func (definition *ProcessDefinition) Fn(processor ProcessorFn) *ProcessorDefinition {
	return definition.P(NewProcessorFromFn(processor))
}

// Ref --
func (definition *ProcessDefinition) Ref(ref string) *ProcessorDefinition {
	registry := definition.context.Registry()
	ifc, err := registry.Lookup(ref)

	if ifc != nil && err == nil {
		if p, ok := ifc.(Processor); ok {
			return definition.P(p)
		}
	}

	return definition.parent
}
