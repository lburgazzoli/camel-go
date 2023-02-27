package processors

import (
	"fmt"

	"github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Reifyable interface {
	Reify(api.Context) (string, error)
}

type Step struct {
	Reifyable
	api.WithOutputs

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

func (s *Step) Reify(ctx api.Context) (string, error) {
	r, ok := s.t.(Reifyable)
	if !ok {
		return "", camelerrors.InternalError("non reifiable step")
	}

	if o, ok := s.t.(api.OutputAware); ok {
		for _, pid := range s.Outputs() {
			o.Next(pid)
		}
	}

	return r.Reify(ctx)
}
