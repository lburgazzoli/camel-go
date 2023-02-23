package model

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Reifybiable interface {
	Reify() (actor.Actor, error)
}

type Step struct {
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

	t := factory()

	if err := node.Content[1].Decode(t); err != nil {
		return errors.Wrapf(err, "unable to decode tag: %s (line: %d, column: %d) ", tag, node.Line, node.Column)
	}

	return nil
}
