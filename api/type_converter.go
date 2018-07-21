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

package api

import (
	"fmt"
	"reflect"
)

// ==========================
//
// Global Converters
//
// ==========================

// TypeConverter --
type TypeConverter func(source interface{}, targetType reflect.Type) (interface{}, error)

// NewConbinedTypeConverter --
func NewConbinedTypeConverter(converter TypeConverter, converters ...TypeConverter) TypeConverter {
	return func(source interface{}, targetType reflect.Type) (interface{}, error) {
		var answer interface{}
		var err error

		answer, err = converter(source, targetType)
		if answer != nil && err == nil {
			return answer, err
		}

		for _, c := range converters {
			answer, err = c(source, targetType)
			if answer != nil && err == nil {
				return answer, err
			}
		}

		return nil, fmt.Errorf("unsupported type conversion (source:%v, target:%v", source, targetType)
	}
}
