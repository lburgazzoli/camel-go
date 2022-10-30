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
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/lburgazzoli/camel-go/api"
)

// ==========================
//
// DefaultExpression
//
// ==========================

// DefaultExpression --
type DefaultExpression struct {
	raw       string
	evaluated string
}

// Raw --
func (e *DefaultExpression) Raw() string {
	return e.raw
}

// Evaluated --
func (e *DefaultExpression) Evaluated() string {
	return e.evaluated
}

// Evaluate --
func (e *DefaultExpression) Evaluate(exchange api.Exchange) (string, error) {

	re, err := regexp.Compile(`\${[^\${}]+}`)
	if err != nil {
		return "", err
	}

	variables := re.FindAllString(e.raw, -1)

	e.evaluated = e.raw
	for _, variable := range variables {
		variable = variable[2 : len(variable)-1]

		switch {

		case strings.HasPrefix(variable, simpleEnvPrefix): // Wants an environment variable
			variable = variable[len(simpleEnvPrefix):]
			e.evaluated = strings.ReplaceAll(e.evaluated, "${"+simpleEnvPrefix+variable+"}", os.Getenv(variable))

		case strings.HasPrefix(variable, simpleHeaderPrefix): // Wants a header
			variable = variable[len(simpleHeaderPrefix):]

			path := strings.Split(variable, ".")

			value, ok := exchange.Headers().Lookup(path[0])
			if !ok {
				return "", fmt.Errorf("header %s not found", path[0])
			}

			for i := 1; i < len(path); i++ {
				value, err = getValue(value, path[i])
				if err != nil {
					return "", fmt.Errorf("invalid path %s on key %s (%s)", variable, path[i], err)
				}
			}
			e.evaluated = strings.ReplaceAll(e.evaluated, "${"+simpleHeaderPrefix+variable+"}", fmt.Sprintf("%v", value))

		case variable == "body": // Wants the body
			valueStr := fmt.Sprintf("%v", exchange.Body())
			e.evaluated = strings.ReplaceAll(e.evaluated, "${"+variable+"}", valueStr)

		case strings.HasPrefix(variable, simpleBodyPrefix): // Wants a field of the body
			variable = variable[len(simpleBodyPrefix):]
			path := strings.Split(variable, ".")

			value := exchange.Body()
			for i := 0; i < len(path); i++ {
				value, err = getValue(value, path[i])
				if err != nil {
					return "", fmt.Errorf("invalid path %s on key %s (%s)", variable, path[i], err)
				}
			}
			e.evaluated = strings.ReplaceAll(e.evaluated, "${"+simpleBodyPrefix+variable+"}", fmt.Sprintf("%v", value))
		}
	}
	return e.evaluated, nil
}

// ==========================
//
//
//
// ==========================

// Simple Language
func Simple(raw string) api.Expression {
	return &DefaultExpression{
		raw: raw,
	}
}

// getValue --
func getValue(object any, key string) (v any, e error) {

	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("failed to convert - %v", r)
		}
	}()

	switch reflect.TypeOf(object).Kind() {

	case reflect.Map:
		objValue := reflect.ValueOf(object)

		value := objValue.MapIndex(reflect.ValueOf(key))
		if value.IsZero() {
			return nil, fmt.Errorf("object %v doesn't have value with key %s", object, key)
		}
		return value.Interface(), nil

	case reflect.Array, reflect.Slice:
		objValue := reflect.ValueOf(object)

		index, err := strconv.Atoi(key)
		if err != nil {
			return nil, fmt.Errorf("expected an int, got %T %s", key, key)
		}

		if index < 0 || index >= objValue.Len() {
			return nil, fmt.Errorf("index %v out of range %v", index, objValue.Len())
		}
		return objValue.Index(index).Interface(), nil

	case reflect.Struct:
		objValue := reflect.ValueOf(object)

		field := objValue.FieldByName(key)
		if field.IsZero() {
			return nil, fmt.Errorf("object %v doesn't have field %s", object, key)
		}
		return field.Interface(), nil

	default:
		return nil, fmt.Errorf("expected a map, array or slice to use the key %s, but got %T", key, object)
	}
}

// ForEach --
func ForEachIn(object any, fn func(key any, value any)) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to convert - %v", r)
		}
	}()

	switch reflect.TypeOf(object).Kind() {

	case reflect.Map:
		objValue := reflect.ValueOf(object)

		for _, key := range objValue.MapKeys() {
			fn(key.Interface(), objValue.MapIndex(key).Interface())
		}

	case reflect.Slice, reflect.Array:
		objValue := reflect.ValueOf(object)

		for i := 0; i < objValue.Len(); i++ {
			fn(i, objValue.Index(i).Interface())
		}

	case reflect.Struct:
		objValue := reflect.ValueOf(object)

		for i := 0; i < objValue.NumField(); i++ {
			fn(objValue.Type().Field(i).Name, objValue.Field(i).Interface())
		}

	default:
		fn(0, object)
	}

	return nil
}

// ==========================
//
// this is where constant
// should be set
//
// ==========================

// Simple Header Prefix
const simpleHeaderPrefix = "header."

// Simple Body Prefix
const simpleBodyPrefix = "body."

// Simple Env Prefix
const simpleEnvPrefix = "env."
