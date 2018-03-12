package camel

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/spf13/cast"
)

// ==========================
//
//
//
// ==========================

func init() {
	AddTypeConverter(&ToIntConverter{})
}

// ==========================
//
// Duration converter
//
// ==========================

// ToDuration --
type ToDuration interface {
	ToDuration() (time.Duration, error)
}

// ==========================
//
// ToDurationConverter
//
// ==========================

// ToDurationConverter --
type ToDurationConverter struct {
	TypeConverter
}

// Convert --
func (converter *ToDurationConverter) Convert(source interface{}, targetType reflect.Type) (interface{}, error) {
	if targetType == reflect.TypeOf(time.Duration(0)) {

		var answer interface{}
		var err error

		sourceType := reflect.TypeOf(source)
		sourceKind := sourceType.Kind()

		if sourceKind == reflect.Struct {
			if v, ok := source.(ToDuration); ok {
				answer, err = v.ToDuration()
			} else {
				err = fmt.Errorf("unable to convert struct:%T to:%v", source, targetType)
			}
		} else {
			answer, err = cast.ToDurationE(source)
		}

		return answer, err
	}

	return nil, errors.New("unsupported")
}
