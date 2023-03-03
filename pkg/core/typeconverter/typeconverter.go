package typeconverter

import (
	"fmt"
	"reflect"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
)

type DefaultTypeConverter struct {
	TypeConverters []camel.TypeConverterFn
}

func (tc *DefaultTypeConverter) Convert(input interface{}, outVal reflect.Value) error {
	// largely inspired by mapstructure type conversion logic

	if input == nil {
		return nil
	}

	var inputVal reflect.Value

	if input != nil {
		inputVal = reflect.ValueOf(input)

		if inputVal.Kind() == reflect.Ptr && inputVal.IsNil() {
			input = nil
		}
	}

	if !inputVal.IsValid() {
		return fmt.Errorf("error decoding from %v to %v: inoput type is invalid", inputVal.Type(), outVal.Type())
	}

	for i := range tc.TypeConverters {
		var err error

		input, err = tc.TypeConverters[i](inputVal.Type(), outVal.Type(), inputVal.Interface())

		if err != nil {
			return fmt.Errorf("error decoding from %v to %v", inputVal.Type(), outVal.Type())
		}
	}

	var err error

	outputKind := tc.getKind(outVal)

	switch outputKind {
	case reflect.Bool:
		// err = d.decodeBool(name, input, outVal)
	case reflect.Interface:
		// err = d.decodeBasic(name, input, outVal)
	case reflect.String:
		dataVal := reflect.Indirect(reflect.ValueOf(input))
		outVal.SetString(dataVal.String())
	case reflect.Int:
		// err = d.decodeInt(name, input, outVal)
	case reflect.Uint:
		// err = d.decodeUint(name, input, outVal)
	case reflect.Float32:
		// err = d.decodeFloat(name, input, outVal)
	case reflect.Struct:
		// err = d.decodeStruct(name, input, outVal)
	case reflect.Map:
		// err = d.decodeMap(name, input, outVal)
	case reflect.Ptr:
		// _, err = d.decodePtr(name, input, outVal)
	case reflect.Slice:
		// err = d.decodeSlice(name, input, outVal)
	case reflect.Array:
		// err = d.decodeArray(name, input, outVal)
	case reflect.Func:
		// err = d.decodeFunc(name, input, outVal)
	default:
		// If we reached this point then we weren't able to decode it
		return fmt.Errorf("unsupported type: %s", outputKind)
	}

	return err
}

func (tc *DefaultTypeConverter) getKind(val reflect.Value) reflect.Kind {
	kind := val.Kind()

	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32
	default:
		return kind
	}
}
