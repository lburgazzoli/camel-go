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
	"reflect"
)

// ==========================
//
// Helper 'const'
//
// ==========================

// EmptyStruct --
var EmptyStruct = struct{}{}

// ==========================
//
// Helper type definition
//
// ==========================

// TypeInt --
var TypeInt = reflect.TypeOf(int(0))

// TypeUInt --
var TypeUInt = reflect.TypeOf(uint(0))

// TypeInt8 --
var TypeInt8 = reflect.TypeOf(int8(0))

// TypeUInt8 --
var TypeUInt8 = reflect.TypeOf(uint8(0))

// TypeInt16 --
var TypeInt16 = reflect.TypeOf(int16(0))

// TypeUInt16 --
var TypeUInt16 = reflect.TypeOf(uint16(0))

// TypeInt32 --
var TypeInt32 = reflect.TypeOf(int32(0))

// TypeUInt32 --
var TypeUInt32 = reflect.TypeOf(uint32(0))

// TypeInt64 --
var TypeInt64 = reflect.TypeOf(int64(0))

// TypeUInt64 --
var TypeUInt64 = reflect.TypeOf(uint64(0))

// TypeString --
var TypeString = reflect.TypeOf("")

// ==========================
//
// Helpers
//
// ==========================

// IsInt --
func IsInt(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Int:
		return true
	case reflect.Uint:
		return true
	case reflect.Int8:
		return true
	case reflect.Uint8:
		return true
	case reflect.Int16:
		return true
	case reflect.Uint16:
		return true
	case reflect.Int32:
		return true
	case reflect.Uint32:
		return true
	case reflect.Int64:
		return true
	case reflect.Uint64:
		return true
	default:
		return false
	}
}
