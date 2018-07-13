package api

// ==========================
//
//
//
// ==========================

// Subscription --
type Subscription interface {
	Cancel()
}

func newdefaultSubcription() *defaultSubcription {
	return &defaultSubcription{
		signal: make(chan bool),
	}
}

type defaultSubcription struct {
	signal chan bool
}

func (subscription *defaultSubcription) Cancel() {
	subscription.signal <- true
}

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
	Subscribe(consumer func(Exchange)) Subscription
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
func (processor *defaultProcessor) Subscribe(consumer func(Exchange)) Subscription {
	subscription := newdefaultSubcription()

	go func() {
		for {
			select {
			case exchange := <-processor.out:
				consumer(exchange)
			case _ = <-subscription.signal:
				return
			default:
			}
		}
	}()

	return subscription
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
func (processor *sourceProcessor) Subscribe(consumer func(Exchange)) Subscription {
	subscription := newdefaultSubcription()

	go func() {
		for {
			select {
			case exchange := <-processor.in:
				consumer(exchange)
			case _ = <-subscription.signal:
				return
			default:
			}
		}
	}()

	return subscription
}

// Parent --
func (processor *sourceProcessor) Parent(parent Processor) Processor {
	parent.Subscribe(func(e Exchange) {
		processor.Publish(e)
	})

	return processor
}
