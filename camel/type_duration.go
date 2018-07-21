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
	"time"

	"github.com/spf13/cast"
)

// ==========================
//
// Init
//
// ==========================

func init() {
	RootContext.AddTypeConverter(ToDurationConverter)
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

// ToDuratioinConverter --
func ToDurationConverter(source interface{}, targetType reflect.Type) (interface{}, error) {
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
