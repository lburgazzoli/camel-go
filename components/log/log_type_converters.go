package log

import (
	"errors"
	"fmt"
	"reflect"

	zl "github.com/rs/zerolog"

	"github.com/lburgazzoli/camel-go/camel"
)

func init() {
	camel.AddTypeConverter(&stringToLogLevelConverter{})

}

// ==========================
//
//
//
// ==========================

type stringToLogLevelConverter struct {
}

func (lc stringToLogLevelConverter) Convert(source interface{}, targetType reflect.Type) (interface{}, error) {

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
		if l, ok := source.(int8); ok {
			switch l {
			case int8(zl.DebugLevel):
				return zl.DebugLevel, nil
			case int8(zl.InfoLevel):
				return zl.InfoLevel, nil
			case int8(zl.WarnLevel):
				return zl.WarnLevel, nil
			case int8(zl.FatalLevel):
				return zl.FatalLevel, nil
			case int8(zl.PanicLevel):
				return zl.PanicLevel, nil
			default:
				return nil, fmt.Errorf("unknown level %d", l)
			}
		}
	}

	return nil, errors.New("unsupported")
}
