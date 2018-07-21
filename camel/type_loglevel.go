// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
