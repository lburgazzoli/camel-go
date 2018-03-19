package camel

import (
	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/iterable"
	"github.com/reactivex/rxgo/observable"
	"github.com/rs/zerolog/log"
)

// ==========================
//
//
//
// ==========================

// NewPipe --
func NewPipe() *Pipe {
	p := Pipe{}
	p.in = make(chan interface{})

	iter, _ := iterable.New(p.in)
	p.observable = observable.From(iter)

	return &p
}

// Pipe --
type Pipe struct {
	in         chan interface{}
	observable observable.Observable
}

// Subscribe --
func (pipe *Pipe) Subscribe(consumer func(*Exchange)) *Pipe {
	onNext := handlers.NextFunc(func(item interface{}) {
		if exchange, ok := item.(*Exchange); ok {
			consumer(exchange)
		} else {
			log.Panic().Msgf("unexpected type: %T", item)
		}
	})

	pipe.observable.Subscribe(onNext)

	return pipe
}

// Publish --
func (pipe *Pipe) Publish(exchange *Exchange) *Pipe {
	pipe.in <- exchange

	return pipe
}

// PublishAsync --
func (pipe *Pipe) PublishAsync(exchange *Exchange) *Pipe {
	go pipe.Publish(exchange)

	return pipe
}

// NewProcessorPipe --
func NewProcessorPipe(parent *Pipe, processor Processor, processors ...Processor) *Pipe {
	next := NewPipe()

	parent.Subscribe(func(e *Exchange) {
		processor.Process(e)

		for _, proc := range processors {
			proc.Process(e)
		}

		next.Publish(e)
	})

	return next
}

// NewTransformerPipe --
func NewTransformerPipe(pipe *Pipe, processor Trasformer, processors ...Trasformer) *Pipe {
	next := NewPipe()

	pipe.Subscribe(func(e *Exchange) {
		e = processor.Trasform(e)

		for _, proc := range processors {
			e = proc.Trasform(e)
		}

		next.Publish(e)
	})

	return next
}

// NewPredicatePipe --
func NewPredicatePipe(pipe *Pipe, processor Predicate, processors ...Predicate) *Pipe {
	next := NewPipe()

	pipe.Subscribe(func(e *Exchange) {
		if ok := processor.Test(e); !ok {
			return
		}

		for _, proc := range processors {
			if ok := proc.Test(e); !ok {
				return
			}
		}

		next.Publish(e)
	})

	return next
}
