package processor

import "github.com/lburgazzoli/camel-go/api"

// Fn --
type Fn func(api.Exchange, chan<- api.Exchange)

// ==========================
//
//
//
// ==========================

// Connect --
func Connect(source api.Processor, destination api.Processor) {
	source.Subscribe(func(exchange api.Exchange) {
		destination.Publish(exchange)
	})
}

// ==========================
//
//
//
// ==========================

// New --
func New(fn Fn) api.Processor {
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

// NewProcessingPipeline --
func NewProcessingPipeline(consumer func(api.Exchange), consumers ...func(api.Exchange)) api.Processor {
	var fn Fn

	if len(consumers) > 0 {
		fn = func(exchange api.Exchange, out chan<- api.Exchange) {
			consumer(exchange)

			for _, c := range consumers {
				c(exchange)
			}

			out <- exchange
		}
	} else {
		fn = func(exchange api.Exchange, out chan<- api.Exchange) {
			consumer(exchange)

			out <- exchange
		}
	}

	return New(fn)
}

// NewFilteringPipeline --
func NewFilteringPipeline(consumer func(api.Exchange) bool, consumers ...func(api.Exchange) bool) api.Processor {
	var fn Fn

	if len(consumers) > 0 {
		fn = func(exchange api.Exchange, out chan<- api.Exchange) {
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
	} else {
		fn = func(exchange api.Exchange, out chan<- api.Exchange) {
			if c := consumer(exchange); c {
				out <- exchange
			}
		}
	}

	return New(fn)
}

// ==========================
//
//
//
// ==========================

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

// defaultProcessor --
type defaultProcessor struct {
	api.Processor

	in  chan api.Exchange
	out chan api.Exchange
	fn  func(api.Exchange, chan<- api.Exchange)
}

// Publish --
func (processor *defaultProcessor) Publish(exchange api.Exchange) {
	processor.in <- exchange
}

// Subscribe --
func (processor *defaultProcessor) Subscribe(consumer func(api.Exchange)) api.Subscription {
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
