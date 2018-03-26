package camel

// ==========================
//
//
//
// ==========================

// Transformer --
type Transformer interface {
	Transform(*Exchange) *Exchange
}

// TransformerFn --
type TransformerFn func(*Exchange) *Exchange

type transformerFnBridge struct {
	Processor
	fn TransformerFn
}

func (bridge *transformerFnBridge) Transform(e *Exchange) *Exchange {
	return bridge.fn(e)
}

// NewTransformerFromFn --
func NewTransformerFromFn(fn TransformerFn) Transformer {
	return &transformerFnBridge{fn: fn}
}
