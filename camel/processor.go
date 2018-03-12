package camel

// ==========================
//
// Processor
//
// ==========================

// Processor --
type Processor interface {
	Process(exchange *Exchange)
}

// ==========================
//
// Pipeline
//
// ==========================

// Pipeline --
type Pipeline struct {
	Processor

	context    *Context
	processors []Processor
}

// NewPipeline --
func NewPipeline(context *Context, processors []Processor) *Pipeline {
	return &Pipeline{
		context:    context,
		processors: processors,
	}
}

// Process --
func (pipeline *Pipeline) Process(exchange *Exchange) {
	for _, processor := range pipeline.processors {
		processor.Process(exchange)
	}
}
