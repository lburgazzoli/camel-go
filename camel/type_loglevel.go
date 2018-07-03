package camel

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/rs/zerolog"
)

// ==========================
//
// Init
//
// ==========================

func init() {
	RootContext.AddTypeConverter(ToLogLevelConverter)
}

// ==========================
//
//
//
// ==========================

// TypeLogLevel --
var TypeLogLevel = reflect.TypeOf(zerolog.InfoLevel)

// ToLogLevelConverter --
func ToLogLevelConverter(source interface{}, targetType reflect.Type) (interface{}, error) {
	if targetType == TypeLogLevel {
		if l, ok := source.(string); ok {
			switch l {
			case "debug":
				return zerolog.DebugLevel, nil
			case "info":
				return zerolog.InfoLevel, nil
			case "warn":
				return zerolog.WarnLevel, nil
			case "fatal":
				return zerolog.FatalLevel, nil
			case "panic":
				return zerolog.PanicLevel, nil
			default:
				return nil, fmt.Errorf("unknown level %s", l)
			}
		}
	}

	return nil, errors.New("unsupported")
}
