package camel

import (
	"log"
	"reflect"
	"strings"
)

// SetProperty --
func SetProperty(context *Context, target interface{}, name string, value interface{}) bool {
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

			log.Fatalf("unable to set field through method call (name=%s, target=%v, error=%v)",
				name,
				target,
				err,
			)
		}
	} else {
		log.Fatalf("unable to set field %s on %v as it is not a pointer", name, target)
	}

	return false
}

// SetProperties --
func SetProperties(context *Context, target interface{}, options map[string]interface{}) int {
	count := 0

	for k, v := range options {
		if SetProperty(context, target, k, v) {
			count++

			// remove the property if successfully
			// se to the target
			delete(options, k)
		}
	}

	return count
}
