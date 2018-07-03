package camel

import "github.com/lburgazzoli/camel-go/api"

// ==========================
//
//
//
// ==========================

// ProcessorFn --
type ProcessorFn func(api.Exchange, chan<- api.Exchange)

// Processor --
type Processor interface {
	Publish(api.Exchange) Processor
	Subscribe(consumer func(api.Exchange)) Processor
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
		in:  make(chan api.Exchange),
		out: make(chan api.Exchange),
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

	parent.Subscribe(func(e api.Exchange) {
		p.Publish(e)
	})

	return p
}

// defaultProcessor --
type defaultProcessor struct {
	in  chan api.Exchange
	out chan api.Exchange
	fn  ProcessorFn
}

// Publish --
func (processor *defaultProcessor) Publish(exchange api.Exchange) Processor {
	processor.in <- exchange

	return processor
}

// Subscribe --
func (processor *defaultProcessor) Subscribe(consumer func(api.Exchange)) Processor {
	go func() {
		for exchange := range processor.out {
			consumer(exchange)
		}
	}()

	return processor
}

// Parent --
func (processor *defaultProcessor) Parent(parent Processor) Processor {
	parent.Subscribe(func(e api.Exchange) {
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
		in: make(chan api.Exchange),
	}

	return &p
}

// NewProcessorSourceWithParent --
func NewProcessorSourceWithParent(parent Processor, fn ProcessorFn) Processor {
	p := NewProcessorSource()

	parent.Subscribe(func(e api.Exchange) {
		p.Publish(e)
	})

	return p
}

// defaultProcessor --
type sourceProcessor struct {
	in chan api.Exchange
}

// Publish --
func (processor *sourceProcessor) Publish(exchange api.Exchange) Processor {
	processor.in <- exchange

	return processor
}

// Subscribe --
func (processor *sourceProcessor) Subscribe(consumer func(api.Exchange)) Processor {
	go func() {
		for exchange := range processor.in {
			consumer(exchange)
		}
	}()

	return processor
}

// Parent --
func (processor *sourceProcessor) Parent(parent Processor) Processor {
	parent.Subscribe(func(e api.Exchange) {
		processor.Publish(e)
	})

	return processor
}
