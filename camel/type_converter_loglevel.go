package camel

import (
	"errors"
	"fmt"
	"reflect"

	zl "github.com/rs/zerolog"
)

// ==========================
//
//
//
// ==========================

// ToLogLevelConverter --
func ToLogLevelConverter(source interface{}, targetType reflect.Type) (interface{}, error) {
	if targetType == reflect.TypeOf(zl.InfoLevel) {
		if l, ok := source.(string); ok {
			switch l {
			case "debug":
				return zl.DebugLevel, nil
			case "info":
				return zl.InfoLevel, nil
			case "warn":
				return zl.WarnLevel, nil
			case "fatal":
				return zl.FatalLevel, nil
			case "panic":
				return zl.PanicLevel, nil
			default:
				return nil, fmt.Errorf("unknown level %s", l)
			}
		}
	}

	return nil, errors.New("unsupported")
}
