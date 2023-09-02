package dsl

import (
	"encoding/json"
	"fmt"

	camelApi "github.com/lburgazzoli/camel-go/api/camel/v2alpha1"
	"gopkg.in/yaml.v3"
)

func ToYamlDSL(flows []camelApi.Flow) ([]byte, error) {
	data, err := json.Marshal(&flows)
	if err != nil {
		return nil, err
	}

	jsondata := make([]map[string]interface{}, 0)

	err = json.Unmarshal(data, &jsondata)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %w", err)
	}

	yamldata, err := yaml.Marshal(&jsondata)
	if err != nil {
		return nil, fmt.Errorf("error marshalling to yaml: %w", err)
	}

	return yamldata, nil
}

func ToUnstructured(flows []camelApi.Flow) ([]map[string]interface{}, error) {
	data, err := json.Marshal(&flows)
	if err != nil {
		return nil, err
	}

	answer := make([]map[string]interface{}, 0)

	if err := json.Unmarshal(data, &answer); err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %w", err)
	}

	return answer, nil
}
