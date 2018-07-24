// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package processor

import (
	"github.com/lburgazzoli/camel-go/api"
)

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
		for {
			select {
			case exchange := <-p.in:
				p.fn(exchange, p.out)
			}
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

// NewProcessingService --
func NewProcessingService(service api.Service, processor api.Processor) api.ProcessingService {
	answer := defaultProcessingService{}
	answer.processor = processor
	answer.startFn = service.Start
	answer.stopFn = service.Stop

	if staged, ok := service.(api.StagedService); ok {
		answer.stageFn = staged.Stage
	} else {
		answer.stageFn = func() api.ServiceStage {
			return api.ServiceStageOther
		}
	}

	return &answer
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
			}
		}
	}()

	return subscription
}

// ==========================
//
//
//
// ==========================

type defaultProcessingService struct {
	startFn   func()
	stopFn    func()
	stageFn   func() api.ServiceStage
	processor api.Processor
}

func (target *defaultProcessingService) Publish(exchange api.Exchange) {
	target.processor.Publish(exchange)
}

func (target *defaultProcessingService) Subscribe(consumer func(api.Exchange)) api.Subscription {
	return target.processor.Subscribe(consumer)
}

func (target *defaultProcessingService) Start() {
	target.startFn()
}

func (target *defaultProcessingService) Stop() {
	target.stopFn()
}

func (target *defaultProcessingService) Stage() api.ServiceStage {
	return target.stageFn()
}
