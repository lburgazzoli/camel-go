package camel

import "reflect"

// SetField --
func SetField(context *Context, target interface{}, field string, value interface{}) {
	v := reflect.ValueOf(target).Elem().FieldByName(field)
	if v.IsValid() {
		converter := context.TypeConverter()
		result, err := converter.Convert(value, v.Type())

		if err != nil {
			newValue := reflect.ValueOf(result)
			v.Set(newValue)
		}
	}
}
