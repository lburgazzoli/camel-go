package route

import (
	"github.com/lburgazzoli/camel-go/pkg/core/model"
	"github.com/lburgazzoli/camel-go/pkg/core/model/endpoint"
)

const TAG = "route"

func init() {
	model.Types[TAG] = func() interface{} {
		return &Route{}
	}
}

type Route struct {
	ID    string            `yaml:"id,omitempty"`
	Group string            `yaml:"group,omitempty"`
	From  endpoint.Endpoint `yaml:"from"`
	Steps []model.Step      `yaml:"steps,omitempty"`
}

func (r *Route) Reify() error {
	return nil
}
