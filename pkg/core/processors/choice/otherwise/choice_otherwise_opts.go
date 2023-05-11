package otherwise

import (
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
)

type OptionFn func(*Otherwise)

func WithStep(step processors.Step) OptionFn {
	return func(in *Otherwise) {
		in.Steps = append(in.Steps, step)
	}
}

func WithProcessor(processor camel.Processor) OptionFn {
	step := processors.NewStep(support.NewProcessorsVerticle(processor))

	return WithStep(step)
}
