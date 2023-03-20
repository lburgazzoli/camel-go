package typeconverter

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/mitchellh/mapstructure"
)

func StringToRawJSON() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {

		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(camel.RawJSON{}) {
			return data, nil
		}

		var result camel.RawJSON

		in, ok := data.(string)
		if !ok {
			return nil, errors.New("type is not []string")
		}

		if err := json.Unmarshal([]byte(in), &result); err != nil {
			return nil, err
		}

		return result, nil
	}
}

func BytesToRawJSON() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {

		if f != reflect.TypeOf([]byte{}) {
			return data, nil
		}
		if t != reflect.TypeOf(camel.RawJSON{}) {
			return data, nil
		}

		var result camel.RawJSON

		in, ok := data.([]byte)
		if !ok {
			return nil, errors.New("type is not []byte")
		}

		if err := json.Unmarshal(in, &result); err != nil {
			return nil, err
		}

		return result, nil
	}
}
