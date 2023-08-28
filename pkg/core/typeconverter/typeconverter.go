package typeconverter

import (
	"reflect"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/mitchellh/mapstructure"
)

type Fn func(reflect.Type, reflect.Type, interface{}) (interface{}, error)

var TypeConverters = make([]Fn, 0)

func NewDefaultTypeConverter() (camel.TypeConverter, error) {

	hooks := make([]mapstructure.DecodeHookFunc, 0, len(TypeConverters))
	for i := range TypeConverters {
		hooks[i] = TypeConverters[i]
	}

	hooks = append(hooks, mapstructure.StringToTimeDurationHookFunc())
	hooks = append(hooks, StringToRawJSON())
	hooks = append(hooks, BytesToRawJSON())
	hooks = append(hooks, TimeToBytes())
	hooks = append(hooks, TimeToString())

	tc := defaultTypeConverter{
		decodeHook: mapstructure.ComposeDecodeHookFunc(hooks...),
	}

	return &tc, nil
}

type defaultTypeConverter struct {
	decodeHook mapstructure.DecodeHookFunc
}

func (tc *defaultTypeConverter) Convert(input interface{}, output interface{}) (bool, error) {
	if input == nil {
		return false, nil
	}

	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           output,
		DecodeHook:       tc.decodeHook,
	})

	if err != nil {
		return false, err
	}

	if err := dec.Decode(input); err != nil {
		return false, err
	}

	return true, nil
}
