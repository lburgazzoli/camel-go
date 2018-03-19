package camel

// ==========================
//
//
//
// ==========================

// Processor --
type Processor interface {
	Process(*Exchange)
}

// ProcessorFn --
type ProcessorFn func(*Exchange)

type processorFnBridge struct {
	Processor
	fn ProcessorFn
}

func (bridge *processorFnBridge) Process(e *Exchange) {
	bridge.fn(e)
}

// NewProcessorFromFn --
func NewProcessorFromFn(fn ProcessorFn) Processor {
	return &processorFnBridge{fn: fn}
}
