package camel

// ==========================
//
//
//
// ==========================

// Trasformer --
type Trasformer interface {
	Trasform(*Exchange) *Exchange
}

// TrasformerFn --
type TrasformerFn func(*Exchange) *Exchange

type trasformerFnBridge struct {
	Processor
	fn TrasformerFn
}

func (bridge *trasformerFnBridge) Trasform(e *Exchange) *Exchange {
	return bridge.fn(e)
}

// NewTrasformerFromFn --
func NewTrasformerFromFn(fn TrasformerFn) Trasformer {
	return &trasformerFnBridge{fn: fn}
}
