package processors

import (
	"context"
	"fmt"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Reifyable interface {
	Reify(context.Context, camel.Context) (string, error)
}

type Step struct {
	Reifyable
	camel.WithOutputs

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

func (s *Step) Reify(ctx context.Context, camelContext camel.Context) (string, error) {
	r, ok := s.t.(Reifyable)
	if !ok {
		return "", camelerrors.InternalError("non reifiable step")
	}

	if o, ok := s.t.(camel.OutputAware); ok {
		for _, pid := range s.Outputs() {
			o.Next(pid)
		}
	}

	return r.Reify(ctx, camelContext)
}

func NewDefaultVerticle() DefaultVerticle {
	return DefaultVerticle{
		Identity: uuid.New(),
	}
}

type DefaultVerticle struct {
	camel.Identifiable
	camel.WithOutputs

	Identity string `yaml:"id"`

	context camel.Context
}

func (v *DefaultVerticle) Context() camel.Context {
	return v.context
}

func (v *DefaultVerticle) SetContext(ctx camel.Context) {
	v.context = ctx
}

func (v *DefaultVerticle) ID() string {
	return v.Identity
}

func (v *DefaultVerticle) Dispatch(msg camel.Message) {

	for _, id := range v.Outputs() {
		if err := v.context.Send(id, msg); err != nil {
			panic(err)
		}
	}
}
