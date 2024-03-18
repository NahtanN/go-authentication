package utils

import "reflect"

func ModelHasColumn(model interface{}, column string) bool {
	modelFields := reflect.TypeOf(model)

	for i := 0; i < modelFields.NumField(); i++ {
		field := modelFields.Field(i)

		if field.Tag.Get("db") == column {
			return true
		}
	}

	return false
}
