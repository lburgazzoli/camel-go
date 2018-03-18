package camel

import "reflect"

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

// Fn --
func (definition *FilterDefinition) Fn(predicate Predicate) *ProcessorDefinition {
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

// Ref --
func (definition *FilterDefinition) Ref(ref string) *ProcessorDefinition {
	registry := definition.context.Registry()
	ifc, err := registry.Lookup(ref)

	if ifc != nil && err == nil {
		if IsPredicate(ifc) {
			p := func(e *Exchange) bool {
				pv := reflect.ValueOf(ifc)
				ev := reflect.ValueOf(e)
				rv := pv.Call([]reflect.Value{ev})

				return rv[0].Bool()
			}

			return definition.Fn(p)
		}
	}

	return definition.child
}
