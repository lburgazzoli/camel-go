package wasm

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/util/registry"
	"github.com/lburgazzoli/camel-go/pkg/wasm"
)

func New(opts ...OptionFn) *Wasm {
	answer := &Wasm{}
	answer.Name = "process"

	for _, o := range opts {
		o(answer)
	}

	return answer
}

type Definition struct {
	Name  string `yaml:"name"`
	Path  string `yaml:"path"`
	Image string `yaml:"image,omitempty"`
}

type Wasm struct {
	Definition `yaml:",inline"`
}

func (l *Wasm) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		return l.UnmarshalText([]byte(value.Value))
	case yaml.MappingNode:
		if err := value.Decode(&l.Definition); err != nil {
			return err
		}

		if l.Name == "" {
			l.Name = "process"
		}

		return nil
	default:
		return fmt.Errorf("unsupported node kind: %v (line: %d, column: %d)", value.Kind, value.Line, value.Column)
	}
}

func (l *Wasm) UnmarshalText(text []byte) error {
	in := string(text)
	parts := strings.Split(in, "?")

	switch len(parts) {
	case 1:
		mn := strings.Split(parts[0], "#")
		switch len(mn) {
		case 1:
			l.Path = mn[0]
		case 2:
			l.Path = mn[0]
			l.Name = mn[1]
		default:
			return camelerrors.InvalidParameterf("wasm", "unsupported wasm reference '%s'", in)
		}
	case 2:
		l.Image = parts[0]

		mn := strings.Split(parts[1], "#")
		switch len(mn) {
		case 1:
			l.Path = mn[0]
		case 2:
			l.Path = mn[0]
			l.Name = mn[1]
		default:
			return camelerrors.InvalidParameterf("wasm", "unsupported wasm reference '%s'", in)
		}
	default:
		return camelerrors.InvalidParameterf("wasm", "unsupported wasm reference '%s'", in)
	}

	if l.Name == "" {
		l.Name = "process"
	}

	return nil
}

func (l *Wasm) Predicate(_ context.Context, _ camel.Context) (camel.Predicate, error) {
	return nil, camelerrors.NotImplemented("NotSupported")
}

func (l *Wasm) Processor(ctx context.Context, _ camel.Context) (camel.Processor, error) {
	if l.Path == "" {
		return nil, camelerrors.MissingParameterf("wasm.path", "failure configuring wasm processor")
	}

	r, err := wasm.NewRuntime(ctx)
	if err != nil {
		return nil, err
	}

	var reader io.ReadCloser

	if l.Image != "" {
		content, err := registry.Blob(ctx, l.Image, l.Path)
		if err != nil {
			return nil, err
		}

		reader = content
	} else {
		content, err := os.Open(l.Path)
		if err != nil {
			return nil, err
		}

		reader = content
	}

	module, err := r.Load(ctx, reader)
	if err != nil {
		return nil, err
	}

	proc, err := module.Processor(ctx, l.Definition.Name)
	if err != nil {
		return nil, err
	}

	p := func(ctx context.Context, m camel.Message) error {
		err := proc.Process(ctx, m)
		if err != nil {
			return err
		}

		return nil
	}

	return p, nil
}

func (l *Wasm) Transformer(_ context.Context, _ camel.Context) (camel.Transformer, error) {
	return nil, camelerrors.NotImplemented("NotSupported")
}
