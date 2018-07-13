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

// simpleSubcription --
type simpleSubcription struct {
	fn func()
}

// Cancel --
func (subscription *simpleSubcription) Cancel() {
	subscription.fn()
}

// ==========================
//
//
//
// ==========================

// Publisher --
type Publisher interface {
	Publish(Exchange)
}

// Subscriber --
type Subscriber interface {
	Subscribe(func(Exchange)) Subscription
}

// Processor --
type Processor interface {
	Publisher
	Subscriber
}

// ==========================
//
//
//
// ==========================

// Connect --
func Connect(source Processor, destination Processor) {
	source.Subscribe(func(exchange Exchange) {
		destination.Publish(exchange)
	})
}

// ==========================
//
//
//
// ==========================

// NewProcessor --
func NewProcessor(fn func(Exchange, chan<- Exchange)) Processor {
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

// NewProcessingPipeline --
func NewProcessingPipeline(consumer func(Exchange), consumers ...func(Exchange)) Processor {
	fn := func(exchange Exchange, out chan<- Exchange) {
		consumer(exchange)

		for _, c := range consumers {
			c(exchange)
		}

		out <- exchange
	}

	return NewProcessor(fn)
}

// NewFilteringPipeline --
func NewFilteringPipeline(consumer func(Exchange) bool, consumers ...func(Exchange) bool) Processor {
	fn := func(exchange Exchange, out chan<- Exchange) {
		if c := consumer(exchange); !c {
			return
		}

		for _, c := range consumers {
			if c := c(exchange); !c {
				return
			}
		}

		out <- exchange
	}

	return NewProcessor(fn)
}

// ==========================
//
//
//
// ==========================

// defaultProcessor --
type defaultProcessor struct {
	Processor

	in  chan Exchange
	out chan Exchange
	fn  func(Exchange, chan<- Exchange)
}

// Publish --
func (processor *defaultProcessor) Publish(exchange Exchange) {
	processor.in <- exchange
}

// Subscribe --
func (processor *defaultProcessor) Subscribe(consumer func(Exchange)) Subscription {
	signal := make(chan bool)
	subscription := &simpleSubcription{
		fn: func() {
			signal <- true
		},
	}

	go func() {
		for {
			select {
			case exchange := <-processor.out:
				consumer(exchange)
			case _ = <-signal:
				return
			default:
			}
		}
	}()

	return subscription
}
