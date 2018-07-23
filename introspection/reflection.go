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

package introspection

import (
	"log"
	"reflect"
	"strings"

	zlog "github.com/rs/zerolog/log"

	"github.com/lburgazzoli/camel-go/api"
)

// SetProperty --
func SetProperty(context api.Context, target interface{}, name string, value interface{}) bool {
	var v reflect.Value
	var t reflect.Value
	var f reflect.Value
	var m reflect.Value

	if reflect.TypeOf(target).Kind() == reflect.Ptr {
		v = reflect.ValueOf(target)
		t = v.Elem()

		// **************************
		// Field
		// **************************

		f = t.FieldByName(name)
		if f.IsValid() && f.CanSet() {
			targetType := f.Type()
			converter := context.TypeConverter()
			result, err := converter(value, targetType)
			if err == nil && result != nil {
				newValue := reflect.ValueOf(result)
				f.Set(newValue)

				return true
			}

			log.Fatalf("unable to set field (name=%s, target=%v, error=%v)",
				name,
				target,
				err,
			)
		}

		// **************************
		// Method
		// **************************

		if !strings.HasPrefix(name, "Set") {
			name = "Set" + strings.ToUpper(name[0:1]) + name[1:]
		}

		m = v.MethodByName(name)
		if m.IsValid() && m.Type().NumIn() == 1 {
			targetType := m.Type().In(0)
			converter := context.TypeConverter()
			result, err := converter(value, targetType)
			if err == nil && result != nil {
				newValue := reflect.ValueOf(result)
				args := []reflect.Value{newValue}

				m.Call(args)

				return true
			}

			zlog.Fatal().Msgf("unable to set field through method call (name=%s, target=%v, error=%v)",
				name,
				target,
				err,
			)
		}
	} else {
		zlog.Fatal().Msgf("unable to set field %s on %v as it is not a pointer", name, target)
	}

	return false
}

// SetProperties --
func SetProperties(context api.Context, target interface{}, options map[string]interface{}) int {
	count := 0

	var k string
	var v interface{}
	var ok bool

	for k, v = range options {
		// check if it is a reference
		if len(k) > 1 && k[0] == '#' {
			k := k[1:]

			// try to lookup value from registry
			if v, ok = context.Registry().Lookup(k); !ok {
				zlog.Fatal().Msgf("unable to find %s from registr", k)
			}
		}

		if SetProperty(context, target, k, v) {
			count++

			// remove the property if successfully
			// se to the target
			delete(options, k)
		}
	}

	return count
}
