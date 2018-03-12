package camel

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

// ==========================
//
// Boolean converter
//
// ==========================

// ToBool --
type ToBool interface {
	ToBool() (bool, error)
}

// ==========================
//
// ToBoolConverter
//
// ==========================

// ToBoolConverter --
func ToBoolConverter(source interface{}, targetType reflect.Type) (interface{}, error) {
	if targetType == reflect.TypeOf(true) {

		var answer interface{}
		var err error

		sourceType := reflect.TypeOf(source)
		sourceKind := sourceType.Kind()

		if sourceKind == reflect.Struct {
			if v, ok := source.(ToBool); ok {
				answer, err = v.ToBool()
			} else {
				err = fmt.Errorf("unable to convert struct:%T to:%v", source, targetType)
			}
		} else {
			answer, err = cast.ToBoolE(source)
		}

		return answer, err
	}

	return nil, errors.New("unsupported")
}
