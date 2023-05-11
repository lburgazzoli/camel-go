package when

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/core/language"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

func New(opts ...OptionFn) *When {
	answer := &When{
		DefaultStepsVerticle: processors.NewDefaultStepsVerticle(),
	}

	for _, o := range opts {
		o(answer)
	}

	return answer
}

type When struct {
	processors.DefaultStepsVerticle `yaml:",inline"`
	language.Language               `yaml:",inline"`

	predicate camel.Predicate
	pid       *actor.PID
}

func (w *When) Reify(ctx context.Context) (camel.Verticle, error) {
	c := camel.ExtractContext(ctx)

	w.DefaultVerticle.SetContext(c)

	switch {
	case w.Jq != nil:
		p, err := w.Jq.Predicate(ctx, c)
		if err != nil {
			return nil, err
		}

		w.predicate = p
	default:
		return nil, camelerrors.MissingParameterf("jq", "failure processing %s", "when")
	}

	return w, nil
}

func (w *When) Matches(ctx context.Context, msg camel.Message) (bool, error) {
	if w.predicate == nil {
		return false, camelerrors.InternalErrorf("not configured")
	}

	return w.predicate(ctx, msg)
}
