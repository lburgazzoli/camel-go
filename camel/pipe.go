package camel

import (
	"github.com/rs/zerolog/log"
)

// NewPipe --
func NewPipe() *Pipe {
	return &Pipe{
		Done: make(chan bool, 1),
		In:   nil,
		Next: nil,
	}
}

// NewPipeIn --
func NewPipeIn() *Pipe {
	return &Pipe{
		Done: make(chan bool, 1),
		In:   make(chan *Exchange),
		Next: nil,
	}
}

// NewPipeWithNext --
func NewPipeWithNext(pipe *Pipe) *Pipe {
	return &Pipe{
		Done: pipe.Done,
		In:   nil,
		Next: pipe,
	}
}

// Pipe --
type Pipe struct {
	Done chan bool
	In   chan *Exchange
	Next *Pipe
}

// Publish --
func (pipe *Pipe) Publish(exchange *Exchange) *Pipe {
	if pipe.Next != nil && pipe.Next.In != nil {
		pipe.Next.In <- exchange
	}

	return pipe
}

// PublishAsync --
func (pipe *Pipe) PublishAsync(exchange *Exchange) *Pipe {
	go pipe.Publish(exchange)

	return pipe
}

// ==========================
//
// Helpers
//
// ==========================

// NewProcessorPipe --
func NewProcessorPipe(processor Processor, processors ...Processor) *Pipe {
	pipe := Pipe{}
	pipe.In = nil
	pipe.Done = nil
	pipe.Next = nil

	go func() {
		for {
			select {
			case exchange, ok := <-pipe.In:
				if !ok {
					log.Warn().Msgf("Channel %+v is not ready", pipe.In)
				} else {
					processor(exchange)

					for _, proc := range processors {
						proc(exchange)
					}

					pipe.Publish(exchange)
				}
			case <-pipe.Done:
				log.Info().Msg("done")
				return
			}
		}
	}()

	return &pipe
}

// NewTrasformerPipe --
func NewTrasformerPipe(processor Trasformer, processors ...Trasformer) *Pipe {
	pipe := Pipe{}
	pipe.In = nil
	pipe.Done = nil
	pipe.Next = nil

	go func() {
		for {
			select {
			case exchange, ok := <-pipe.In:
				if !ok {
					log.Warn().Msgf("Channel %+v is not ready", pipe.In)
				} else {
					exchange = processor(exchange)

					for _, proc := range processors {
						exchange = proc(exchange)
					}

					pipe.Publish(exchange)
				}
			case <-pipe.Done:
				log.Info().Msg("done")
				return
			}
		}
	}()

	return &pipe
}
