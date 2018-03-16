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

// Next --
func (pipe *Pipe) Next(next *Pipe) *Pipe {
	return pipe.Subscribe(func(e *Exchange) {
		next.Publish(e)
	})
}

// Subscribe --
func (pipe *Pipe) Subscribe(processor Processor) *Pipe {
	onNext := handlers.NextFunc(func(item interface{}) {
		if exchange, ok := item.(*Exchange); ok {
			processor(exchange)
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

// Process --
func (pipe *Pipe) Process(processor Processor, processors ...Processor) *Pipe {
	next := NewPipe()

	pipe.Subscribe(func(e *Exchange) {
		processor(e)

		for _, proc := range processors {
			proc(e)
		}

		next.Publish(e)
	})

	return next
}

// Transformer --
func (pipe *Pipe) Transformer(processor Trasformer, processors ...Trasformer) *Pipe {
	next := NewPipe()

	pipe.Subscribe(func(e *Exchange) {
		e = processor(e)

		for _, proc := range processors {
			e = proc(e)
		}

		next.Publish(e)
	})

	return next
}
