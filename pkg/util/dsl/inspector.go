package dsl

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	camelApi "github.com/lburgazzoli/camel-go/api/camel/v2alpha1"
)

func NewInspector() *Inspector {
	return &Inspector{}
}

type Inspector struct {
}

// Extract --.
func (i *Inspector) Extract(flows []camelApi.Flow) (*Metadata, error) {
	definitions, err := ToUnstructured(flows)
	if err != nil {
		return nil, err
	}

	meta := NewMetadata()

	for _, definition := range definitions {
		for k, v := range definition {
			if err := i.parseStep(k, v, meta); err != nil {
				return nil, err
			}
		}
	}

	return meta, nil
}

func (i *Inspector) parseStep(key string, content interface{}, meta *Metadata) error {
	var maybeURI string

	switch t := content.(type) {
	case string:
		maybeURI = t
	case map[string]interface{}:
		for k, v := range t {
			switch k {
			case "steps":
				if steps, stepsFormatOk := v.([]interface{}); stepsFormatOk {
					if err := i.parseStepsParam(steps, meta); err != nil {
						return err
					}
				}
			case "uri":
				if vv, isString := v.(string); isString {
					builtURI := vv
					// Inject parameters into URIs to allow other parts of the operator to inspect them
					if params, pok := t["parameters"]; pok {
						if paramMap, pmok := params.(map[interface{}]interface{}); pmok {
							params := make(map[string]string, len(paramMap))
							for k, v := range paramMap {
								ks := fmt.Sprintf("%v", k)
								vs := fmt.Sprintf("%v", v)
								params[ks] = vs
							}
							builtURI = i.appendParameters(builtURI, params)
						}
					}
					maybeURI = builtURI
				}
			default:
				if _, ok := v.(map[interface{}]interface{}); ok {
					if err := i.parseStep(k, v, meta); err != nil {
						return err
					}
				} else if _, ok := v.(map[string]interface{}); ok {
					if err := i.parseStep(k, v, meta); err != nil {
						return err
					}
				} else if ls, ok := v.([]interface{}); ok {
					for _, el := range ls {
						if err := i.parseStep(k, el, meta); err != nil {
							return err
						}
					}
				}
			}
		}
	case map[interface{}]interface{}:
		for k, v := range t {
			switch k {
			case "steps":
				if steps, stepsFormatOk := v.([]interface{}); stepsFormatOk {
					if err := i.parseStepsParam(steps, meta); err != nil {
						return err
					}
				}
			case "uri":
				if vv, isString := v.(string); isString {
					builtURI := vv
					// Inject parameters into URIs to allow other parts of the operator to inspect them
					if params, pok := t["parameters"]; pok {
						if paramMap, pmok := params.(map[interface{}]interface{}); pmok {
							params := make(map[string]string, len(paramMap))
							for k, v := range paramMap {
								ks := fmt.Sprintf("%v", k)
								vs := fmt.Sprintf("%v", v)
								params[ks] = vs
							}
							builtURI = i.appendParameters(builtURI, params)
						}
					}
					maybeURI = builtURI
				}
			default:
				// Always follow children because from/to uris can be nested
				if ks, ok := k.(string); ok {
					if _, ok := v.(map[interface{}]interface{}); ok {
						if err := i.parseStep(ks, v, meta); err != nil {
							return err
						}
					} else if _, ok := v.(map[string]interface{}); ok {
						if err := i.parseStep(ks, v, meta); err != nil {
							return err
						}
					} else if ls, ok := v.([]interface{}); ok {
						for _, el := range ls {
							if err := i.parseStep(ks, el, meta); err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}

	if maybeURI != "" {
		switch key {
		case "from":
			meta.FromURIs = append(meta.FromURIs, maybeURI)
		case "to", "to-d", "toD", "wire-tap", "wireTap":
			meta.ToURIs = append(meta.ToURIs, maybeURI)
		}
	}

	return nil
}

// TODO nolint: gocyclo.
func (i *Inspector) parseStepsParam(steps []interface{}, meta *Metadata) error {
	for _, raw := range steps {
		if step, stepFormatOk := raw.(map[interface{}]interface{}); stepFormatOk {

			if len(step) != 1 {
				return fmt.Errorf("unable to parse step: %v", step)
			}

			for k, v := range step {
				switch kt := k.(type) {
				case fmt.Stringer:
					if err := i.parseStep(kt.String(), v, meta); err != nil {
						return err
					}
				case string:
					if err := i.parseStep(kt, v, meta); err != nil {
						return err
					}
				default:
					return fmt.Errorf("unknown key type: %v, step: %v", k, step)
				}
			}
		}
	}
	return nil
}

func (i *Inspector) appendParameters(uri string, params map[string]string) string {
	prefix := "&"
	if !strings.Contains(uri, "?") {
		prefix = "?"
	}
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		uri += fmt.Sprintf("%s%s=%s", prefix, url.QueryEscape(k), url.QueryEscape(params[k]))
		prefix = "&"
	}
	return uri
}
