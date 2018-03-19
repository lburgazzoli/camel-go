package camel

// ==========================
//
// FilterDefinition
//
//    WORK IN PROGRESS
//
// ==========================

// Filter --
func (definition *ProcessorDefinition) Filter() *FilterDefinition {
	filter := FilterDefinition{}
	filter.context = definition.context
	filter.parent = definition

	definition.child = &filter.ProcessorDefinition

	return &filter
}

// ==========================
//
//
//
//
//
// ==========================

// FilterDefinition --
type FilterDefinition struct {
	ProcessorDefinition
	predicate Predicate
}

// P --
func (definition *FilterDefinition) P(predicate Predicate) *ProcessorDefinition {
	definition.predicate = predicate
	definition.child = NewProcessorDefinitionWithParent(&definition.ProcessorDefinition)

	definition.AddFactory(func(parent *Pipe) (*Pipe, Service) {
		return NewPredicatePipe(parent, predicate), nil
	})

	return definition.child
}

// Fn --
func (definition *FilterDefinition) Fn(predicate PredicateFn) *ProcessorDefinition {
	return definition.P(NewPredicateFromFn(predicate))
}

// Ref --
func (definition *FilterDefinition) Ref(ref string) *ProcessorDefinition {
	registry := definition.context.Registry()
	ifc, err := registry.Lookup(ref)

	if ifc != nil && err == nil {
		if p, ok := ifc.(Predicate); ok {
			return definition.P(p)
		}
	}

	return definition.parent
}
