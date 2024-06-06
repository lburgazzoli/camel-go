package processors

import (
	"context"
	"fmt"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Reifyable interface {
	Reify(context.Context) (camel.Verticle, error)
}

func NewStep(r Reifyable) Step {
	return Step{
		t: r,
	}
}

func ReifySteps(ctx context.Context, steps []Step) ([]camel.Verticle, error) {
	verticles := make([]camel.Verticle, len(steps))

	for s := range steps {
		step := steps[s]

		v, err := step.Reify(ctx)
		if err != nil {
			return nil, err
		}

		verticles[s] = v
	}

	return verticles, nil
}

type Step struct {
	Reifyable

	t interface{}
}

func (s *Step) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("unsupported type (line: %d, column: %d) ", node.Line, node.Column)
	}

	if len(node.Content) != 2 {
		return fmt.Errorf("unsupported content (size: %d)", len(node.Content))
	}

	tag := node.Content[0].Value

	factory, ok := Types[tag]
	if !ok {
		return fmt.Errorf("unsupported tag: %s", tag)
	}

	s.t = factory()

	if err := node.Content[1].Decode(s.t); err != nil {
		return errors.Wrapf(err, "unable to decode tag: %s (line: %d, column: %d) ", tag, node.Line, node.Column)
	}

	return nil
}

func (s *Step) Reify(ctx context.Context) (camel.Verticle, error) {
	r, ok := s.t.(Reifyable)
	if !ok {
		return nil, camelerrors.InternalError("non reifiable step")
	}

	return r.Reify(ctx)
}
