package camel

// ==========================
//
// FilterDefinition
//
//    WORK IN PROGRESS
//
// ==========================

// FilterDefinition --
type FilterDefinition struct {
	ProcessorDefinition
	predicate Predicate
}

// WithPredicate --
func (definition *FilterDefinition) WithPredicate(predicate Predicate) *ProcessorDefinition {
	definition.predicate = predicate
	definition.child = NewProcessorDefinitionWithParent(&definition.ProcessorDefinition)

	definition.addFactory(func(parent *Pipe) (*Pipe, Service) {
		next := NewPipe()

		parent.Subscribe(func(e *Exchange) {
			if definition.predicate != nil && definition.predicate(e) {
				next.Publish(e)
			}
		})

		return next, nil
	})

	return definition.child
}
