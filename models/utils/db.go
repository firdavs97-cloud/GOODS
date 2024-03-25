package utils

import (
	"reflect"
)

func GetColumnsAndValues(object interface{}, skipTags ...string) ([]string, []interface{}) {
	objectType := reflect.TypeOf(object)
	values := make([]interface{}, 0)
	columns := make([]string, 0)
	for i := 0; i < objectType.NumField(); i++ {
		field := objectType.Field(i)
		fieldValue := reflect.ValueOf(object).Field(i)
		skip := false
		for _, s := range skipTags {
			if dbTag := field.Tag.Get(s); dbTag != "" {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		if dbTag := field.Tag.Get("db"); dbTag != "" {
			columns = append(columns, dbTag)
			values = append(values, fieldValue.Interface())
		}
	}
	return columns, values
}
