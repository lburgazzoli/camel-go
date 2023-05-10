package from

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v3"

	camel "github.com/lburgazzoli/camel-go/pkg/api"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors/endpoint"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/pkg/errors"
)

const TAG = "to"

func init() {
	processors.Types[TAG] = func() interface{} {
		return New()
	}
}

func New() *To {
	return &To{
		Definition: Definition{
			Endpoint: endpoint.Endpoint{
				Identity: uuid.New(),
			},
		},
	}
}

type Definition struct {
	endpoint.Endpoint `yaml:",inline"`
}

type To struct {
	Definition `yaml:",inline"`
}

func (t *To) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		return t.UnmarshalText([]byte(value.Value))
	case yaml.MappingNode:
		return value.Decode(&t.Definition)
	default:
		return fmt.Errorf("unsupported node kind: %v (line: %d, column: %d)", value.Kind, value.Line, value.Column)
	}
}

func (t *To) UnmarshalText(text []byte) error {
	t.Endpoint.URI = string(text)
	return nil
}

func (t *To) Reify(ctx context.Context) (camel.Verticle, error) {
	camelContext := camel.ExtractContext(ctx)

	producer, err := t.Endpoint.Producer(camelContext)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating consumer")
	}

	return producer, nil
}
