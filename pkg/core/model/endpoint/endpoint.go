package endpoint

import (
	"github.com/lburgazzoli/camel-go/pkg/core/model"
)

const TAG = "endpoint"

func init() {
	model.Types[TAG] = func() interface{} {
		return &Endpoint{}
	}
}

type Endpoint struct {
	URL        string                 `yaml:"url,omitempty"`
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
}

func (e *Endpoint) Reify() error {
	return nil
}
