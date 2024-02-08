package typeconverter

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/pkg/errors"

	"github.com/go-viper/mapstructure/v2"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
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

func TimeToBytes() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {

		if f != reflect.TypeOf(time.Time{}) {
			return data, nil
		}
		if t != reflect.TypeOf([]byte{}) {
			return data, nil
		}

		in, ok := data.(time.Time)
		if !ok {
			return nil, errors.New("type is not time.Time")
		}

		res := in.Format(time.RFC3339Nano)

		return []byte(res), nil
	}
}

func TimeToString() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {

		if f != reflect.TypeOf(time.Time{}) {
			return data, nil
		}
		if t.Kind() != reflect.String {
			return data, nil
		}

		in, ok := data.(time.Time)
		if !ok {
			return nil, errors.New("type is not time.Time")
		}

		return in.Format(time.RFC3339Nano), nil
	}
}
