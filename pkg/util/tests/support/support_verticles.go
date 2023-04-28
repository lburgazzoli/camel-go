package support

import (
	"context"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/core/verticles"
	"github.com/pkg/errors"
)

type ReifyableVerticle interface {
	camel.Verticle
	processors.Reifyable
}

func NewChannelVerticle(channel chan camel.Message) ReifyableVerticle {
	return &ChannelVerticle{
		DefaultVerticle: processors.NewDefaultVerticle(),
		channel:         channel,
	}
}

type ChannelVerticle struct {
	processors.DefaultVerticle
	channel chan camel.Message
}

func (p *ChannelVerticle) Receive(ac actor.Context) {
	msg, ok := ac.Message().(camel.Message)
	if ok {
		p.channel <- msg

		if ac.Sender() != nil {
			ac.Request(ac.Sender(), msg)
		}
	}
}

func (p *ChannelVerticle) Reify(_ context.Context) (camel.Verticle, error) {
	return p, nil
}

func NewProcessorsVerticle(processor camel.Processor) ReifyableVerticle {
	return &ProcessorVerticle{
		DefaultVerticle: processors.NewDefaultVerticle(),
		processor:       processor,
	}
}

type ProcessorVerticle struct {
	processors.DefaultVerticle
	processor camel.Processor
}

func (p *ProcessorVerticle) Reify(_ context.Context) (camel.Verticle, error) {
	return p, nil
}

func (p *ProcessorVerticle) Receive(ac actor.Context) {
	msg, ok := ac.Message().(camel.Message)
	if ok {
		if err := p.processor(context.Background(), msg); err != nil {
			panic(err)
		}

		if ac.Sender() != nil {
			ac.Request(ac.Sender(), msg)
		}
	}
}

func NewRootVerticle(v camel.Verticle) *RootVerticle {
	return &RootVerticle{
		DefaultVerticle: processors.NewDefaultVerticle(),
		V:               v,
		C:               make(chan camel.Message, 1),
	}
}

type RootVerticle struct {
	processors.DefaultVerticle

	V camel.Verticle
	P *actor.PID
	C chan camel.Message
}

func (r *RootVerticle) Receive(ac actor.Context) {
	switch msg := ac.Message().(type) {
	case *actor.Started:
		p, err := verticles.Spawn(ac, r.V)
		if err != nil {
			panic(errors.Wrapf(err, "unable to spawn verticle with id %s", r.V.ID()))
		}

		r.P = p
	case camel.Message:
		if ac.Sender() != nil {
			for _, c := range ac.Children() {
				if ac.Sender().Equal(c) {
					r.C <- msg
					return
				}
			}
		}

		ac.Request(r.P, msg)
	}
}

func (r *RootVerticle) Get(timeout time.Duration) (camel.Message, error) {
	select {
	case msg := <-r.C:
		return msg, nil
	case <-time.After(timeout):
		return nil, errors.New("timeout")
	}
}
