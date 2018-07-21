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
	RootContext.AddTypeConverter(ToIntConverter)
}

// Integer --
type Integer int

// ==========================
//
// Int converter
//
// ==========================

// ToInt --
type ToInt interface {
	ToInt() (int, error)
}

// ToUInt --
type ToUInt interface {
	ToUInt() (int, error)
}

// ToInt8 --
type ToInt8 interface {
	ToInt8() (int8, error)
}

// ToUInt8 --
type ToUInt8 interface {
	ToUInt8() (int8, error)
}

// ToInt16 --
type ToInt16 interface {
	ToInt16() (int16, error)
}

// ToUInt16 --
type ToUInt16 interface {
	ToUInt16() (int16, error)
}

// ToInt32 --
type ToInt32 interface {
	ToInt32() (int32, error)
}

// ToUInt32 --
type ToUInt32 interface {
	ToUInt32() (int32, error)
}

// ToInt64 --
type ToInt64 interface {
	ToInt64() (int64, error)
}

// ToUInt64 --
type ToUInt64 interface {
	ToUInt64() (int64, error)
}

// ==========================
//
//
//
// ==========================

// ToInt ..
func (target Integer) ToInt() int {
	return int(target)
}

// ToUInt ..
func (target Integer) ToUInt() uint {
	return uint(target)
}

// ToInt8 ..
func (target Integer) ToInt8() int8 {
	return int8(target)
}

// ToUInt8 ..
func (target Integer) ToUInt8() uint8 {
	return uint8(target)
}

// ToInt16 ..
func (target Integer) ToInt16() int16 {
	return int16(target)
}

// ToUInt16 ..
func (target Integer) ToUInt16() uint16 {
	return uint16(target)
}

// ToInt32 ..
func (target Integer) ToInt32() int32 {
	return int32(target)
}

// ToUInt32 ..
func (target Integer) ToUInt32() uint32 {
	return uint32(target)
}

// ToInt64 ..
func (target Integer) ToInt64() int64 {
	return int64(target)
}

// ToUInt64 ..
func (target Integer) ToUInt64() uint64 {
	return uint64(target)
}

// ==========================
//
// ToIntConverter
//
// ==========================

// ToIntConverter --
func ToIntConverter(source interface{}, targetType reflect.Type) (interface{}, error) {
	if !IsInt(targetType) {
		return nil, errors.New("unsupported")
	}

	var answer interface{}
	var err error

	sourceType := reflect.TypeOf(source)
	sourceKind := sourceType.Kind()

	switch targetType.Kind() {
	case reflect.Int:
		if sourceKind == reflect.Struct {
			if v, ok := source.(ToInt); ok {
				answer, err = v.ToInt()
			} else {
				err = fmt.Errorf("unable to convert struct:%T to:%v", source, targetType)
			}
		} else {
			answer, err = cast.ToIntE(source)
		}
	case reflect.Uint:
		if sourceKind == reflect.Struct {
			if v, ok := source.(ToUInt); ok {
				answer, err = v.ToUInt()
			} else {
				err = fmt.Errorf("unable to convert struct:%T to:%v", source, targetType)
			}
		} else {
			answer, err = cast.ToUintE(source)
		}
	case reflect.Int8:
		if sourceKind == reflect.Struct {
			if v, ok := source.(ToInt8); ok {
				answer, err = v.ToInt8()
			} else {
				err = fmt.Errorf("unable to convert struct:%T to:%v", source, targetType)
			}
		} else {
			answer, err = cast.ToInt8E(source)
		}
	case reflect.Uint8:
		if sourceKind == reflect.Struct {
			if v, ok := source.(ToUInt8); ok {
				answer, err = v.ToUInt8()
			} else {
				err = fmt.Errorf("unable to convert struct:%T to:%v", source, targetType)
			}
		} else {
			answer, err = cast.ToUint8E(source)
		}
	case reflect.Int16:
		if sourceKind == reflect.Struct {
			if v, ok := source.(ToInt16); ok {
				answer, err = v.ToInt16()
			} else {
				err = fmt.Errorf("unable to convert struct:%T to:%v", source, targetType)
			}
		} else {
			answer, err = cast.ToInt16E(source)
		}
	case reflect.Uint16:
		if sourceKind == reflect.Struct {
			if v, ok := source.(ToUInt16); ok {
				answer, err = v.ToUInt16()
			} else {
				err = fmt.Errorf("unable to convert struct:%T to:%v", source, targetType)
			}
		} else {
			answer, err = cast.ToUint16E(source)
		}
	case reflect.Int32:
		if sourceKind == reflect.Struct {
			if v, ok := source.(ToInt32); ok {
				answer, err = v.ToInt32()
			} else {
				err = fmt.Errorf("unable to convert struct:%T to:%v", source, targetType)
			}
		} else {
			answer, err = cast.ToInt32E(source)
		}
	case reflect.Uint32:
		if sourceKind == reflect.Struct {
			if v, ok := source.(ToUInt32); ok {
				answer, err = v.ToUInt32()
			} else {
				err = fmt.Errorf("unable to convert struct:%T to:%v", source, targetType)
			}
		} else {
			answer, err = cast.ToUint32E(source)
		}
	case reflect.Int64:
		if sourceKind == reflect.Struct {
			if v, ok := source.(ToInt64); ok {
				answer, err = v.ToInt64()
			} else {
				err = fmt.Errorf("unable to convert struct:%T to:%v", source, targetType)
			}
		} else {
			answer, err = cast.ToInt64E(source)
		}
	case reflect.Uint64:
		if sourceKind == reflect.Struct {
			if v, ok := source.(ToUInt64); ok {
				answer, err = v.ToUInt64()
			} else {
				err = fmt.Errorf("unable to convert struct:%T to:%v", source, targetType)
			}
		} else {
			answer, err = cast.ToUint64E(source)
		}
	}

	return answer, err
}
