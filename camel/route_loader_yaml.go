package camel

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// LoadRouteFromYAMLFile --
func LoadRouteFromYAMLFile(path string) (Definition, error) {
	data, err := ioutil.ReadFile(path)
	if err == nil {
		return nil, err
	}

	return LoadRouteFromYAML(data)
}

// LoadRouteFromYAML --
func LoadRouteFromYAML(data []byte) (Definition, error) {
	flow := Flow{}
	err := yaml.Unmarshal([]byte(data), &flow)

	if err != nil {
		return nil, err
	}

	var route *RouteDefinition

	if err == nil {
		for i, s := range flow.Steps {
			if i == 0 {
				route = From(s.URI)
			} else {
				route = route.To(s.URI)
			}
		}
	}

	return route, nil
}
