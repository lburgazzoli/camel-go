package camel

// ==========================
//
//
//
// ==========================

// Predicate --
type Predicate interface {
	Test(*Exchange) bool
}

// PredicateFn --
type PredicateFn func(*Exchange) bool

type predicateFnBridge struct {
	Processor
	fn PredicateFn
}

func (bridge *predicateFnBridge) Test(e *Exchange) bool {
	return bridge.fn(e)
}

// NewPredicateFromFn --
func NewPredicateFromFn(fn PredicateFn) Predicate {
	return &predicateFnBridge{fn: fn}
}
