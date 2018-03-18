package camel

import (
	"reflect"
)

// ==========================
//
// ProcessDefinition
//
//    WORK IN PROGRESS
//
// ==========================

// ProcessDefinition --
type ProcessDefinition struct {
	ProcessorDefinition
	processor Processor
}

// Fn --
func (definition *ProcessDefinition) Fn(processor Processor) *ProcessorDefinition {
	definition.processor = processor
	definition.child = NewProcessorDefinitionWithParent(&definition.ProcessorDefinition)

	definition.addFactory(func(parent *Pipe) (*Pipe, Service) {
		next := NewPipe()

		parent.Subscribe(func(e *Exchange) {
			processor(e)
			next.Publish(e)
		})

		return next, nil
	})

	return definition.child
}

// Ref --
func (definition *ProcessDefinition) Ref(ref string) *ProcessorDefinition {
	registry := definition.context.Registry()
	ifc, err := registry.Lookup(ref)

	if ifc != nil && err == nil {
		if IsProcessor(ifc) {
			p := func(e *Exchange) {
				pv := reflect.ValueOf(ifc)
				ev := reflect.ValueOf(e)

				pv.Call([]reflect.Value{ev})
			}

			return definition.Fn(p)
		}
	}

	return definition.child
}
