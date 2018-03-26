package camel

// ==========================
//
// FilterDefinition
//
//    WORK IN PROGRESS
//
// ==========================

// Filter --
func (definition *RouteDefinition) Filter() *FilterDefinition {
	filter := FilterDefinition{}
	filter.parent = definition

	definition.child = &filter.RouteDefinition

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
	RouteDefinition
}

// P --
func (definition *FilterDefinition) P(predicate Predicate) *RouteDefinition {
	definition.parent.AddFactory(func(context *Context, parent *Subject) (*Subject, Service, error) {
		return NewSubject().SubscribeWithPredicate(parent, predicate), nil, nil
	})

	return definition.parent
}

// Fn --
func (definition *FilterDefinition) Fn(predicate PredicateFn) *RouteDefinition {
	return definition.P(NewPredicateFromFn(predicate))
}

// Ref --
func (definition *FilterDefinition) Ref(ref string) *RouteDefinition {
	definition.parent.AddFactory(func(context *Context, parent *Subject) (*Subject, Service, error) {
		registry := context.Registry()
		ifc, err := registry.Lookup(ref)

		if ifc != nil && err == nil {
			if p, ok := ifc.(Predicate); ok {
				return NewSubject().SubscribeWithPredicate(parent, p), nil, nil
			}
		}

		// TODO: error handling
		return nil, nil, nil
	})

	return definition.parent
}
