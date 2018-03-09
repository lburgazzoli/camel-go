package camel

import (
	"log"
	"reflect"
	"strings"
)

// SetField --
func SetField(context *Context, target interface{}, name string, value interface{}) {
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
			valueType := reflect.TypeOf(value)
			converter := context.TypeConverter()
			result, err := converter.Convert(value, valueType)
			if err != nil {
				newValue := reflect.ValueOf(result)
				v.Set(newValue)
			} else {
				log.Fatalf("unable to set field (name=%s, target=%v, error=%v)",
					name,
					target,
					err,
				)
			}

			return
		}

		// **************************
		// Method
		// **************************

		if !strings.HasPrefix(name, "Set") {
			name = "Set" + name
		}

		m = v.MethodByName(name)
		if m.IsValid() && m.Type().NumIn() == 1 {
			valueType := reflect.TypeOf(value)
			converter := context.TypeConverter()
			result, err := converter.Convert(value, valueType)
			if err != nil {
				newValue := reflect.ValueOf(result)
				args := []reflect.Value{newValue}

				m.Call(args)
			} else {
				log.Fatalf("unable to set field thorugh method call (name=%s, target=%v, error=%v)",
					name,
					target,
					err,
				)
			}

			return
		}
	} else {
		log.Fatalf("unable to set field %s on %v as it is not a pointer", name, target)
	}
}
