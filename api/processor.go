package api

// ==========================
//
//
//
// ==========================

// ProcessorFn --
type ProcessorFn func(Exchange, chan<- Exchange)

// Processor --
type Processor interface {
	Publish(Exchange) Processor
	Subscribe(consumer func(Exchange)) Processor
	Parent(parent Processor) Processor
}

// ==========================
//
//
//
// ==========================

// NewProcessor --
func NewProcessor(fn ProcessorFn) Processor {
	p := defaultProcessor{
		in:  make(chan Exchange),
		out: make(chan Exchange),
		fn:  fn,
	}

	go func() {
		for exchange := range p.in {
			p.fn(exchange, p.out)
		}
	}()

	return &p
}

// NewProcessorWithParent --
func NewProcessorWithParent(parent Processor, fn ProcessorFn) Processor {
	p := NewProcessor(fn)

	parent.Subscribe(func(e Exchange) {
		p.Publish(e)
	})

	return p
}

// defaultProcessor --
type defaultProcessor struct {
	in  chan Exchange
	out chan Exchange
	fn  ProcessorFn
}

// Publish --
func (processor *defaultProcessor) Publish(exchange Exchange) Processor {
	processor.in <- exchange

	return processor
}

// Subscribe --
func (processor *defaultProcessor) Subscribe(consumer func(Exchange)) Processor {
	go func() {
		for exchange := range processor.out {
			consumer(exchange)
		}
	}()

	return processor
}

// Parent --
func (processor *defaultProcessor) Parent(parent Processor) Processor {
	parent.Subscribe(func(e Exchange) {
		processor.Publish(e)
	})

	return processor
}

// ==========================
//
//
//
// ==========================

// NewProcessorSource --
func NewProcessorSource() Processor {
	p := sourceProcessor{
		in: make(chan Exchange),
	}

	return &p
}

// NewProcessorSourceWithParent --
func NewProcessorSourceWithParent(parent Processor, fn ProcessorFn) Processor {
	p := NewProcessorSource()

	parent.Subscribe(func(e Exchange) {
		p.Publish(e)
	})

	return p
}

// defaultProcessor --
type sourceProcessor struct {
	in chan Exchange
}

// Publish --
func (processor *sourceProcessor) Publish(exchange Exchange) Processor {
	processor.in <- exchange

	return processor
}

// Subscribe --
func (processor *sourceProcessor) Subscribe(consumer func(Exchange)) Processor {
	go func() {
		for exchange := range processor.in {
			consumer(exchange)
		}
	}()

	return processor
}

// Parent --
func (processor *sourceProcessor) Parent(parent Processor) Processor {
	parent.Subscribe(func(e Exchange) {
		processor.Publish(e)
	})

	return processor
}
