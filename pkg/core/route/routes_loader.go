package route

import (
	"io"

	"github.com/lburgazzoli/camel-go/pkg/core/processors/route"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func Load(in io.Reader) ([]route.Route, error) {
	data, err := io.ReadAll(in)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse routes")
	}

	var holder []struct {
		R route.Route `yaml:"route"`
	}

	if err := yaml.Unmarshal(data, &holder); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal routes")
	}

	routes := make([]route.Route, 0, len(holder))
	for i := range holder {
		r := holder[i].R

		if r.ID == "" {
			r.ID = uuid.New()
		}

		routes = append(routes, r)
	}

	return routes, err
}
