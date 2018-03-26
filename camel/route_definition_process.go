package camel

// ==========================
//
// ProcessDefinition
//
//    WORK IN PROGRESS
//
// ==========================

// Process --
func (definition *RouteDefinition) Process() *ProcessDefinition {
	process := ProcessDefinition{}
	process.parent = definition

	definition.child = &process.RouteDefinition

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
	RouteDefinition
}

// P --
func (definition *ProcessDefinition) P(processor Processor) *RouteDefinition {
	definition.parent.AddFactory(func(context *Context, parent *Subject) (*Subject, Service, error) {
		return NewSubject().SubscribeWithProcessor(parent, processor), nil, nil
	})

	return definition.parent
}

// Fn --
func (definition *ProcessDefinition) Fn(processor ProcessorFn) *RouteDefinition {
	return definition.P(NewProcessorFromFn(processor))
}

// Ref --
func (definition *ProcessDefinition) Ref(ref string) *RouteDefinition {
	definition.parent.AddFactory(func(context *Context, parent *Subject) (*Subject, Service, error) {
		registry := context.Registry()
		ifc, err := registry.Lookup(ref)

		if ifc != nil && err == nil {
			if p, ok := ifc.(Processor); ok {
				return NewSubject().SubscribeWithProcessor(parent, p), nil, nil
			}
		}

		// TODO: error handling
		return nil, nil, nil
	})

	return definition.parent
}
