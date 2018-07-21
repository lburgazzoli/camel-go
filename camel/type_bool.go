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

	"github.com/spf13/cast"
)

// ==========================
//
// Init
//
// ==========================

func init() {
	RootContext.AddTypeConverter(ToBoolConverter)
}

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
