package utils

import "reflect"

func ModelHasColumn(model interface{}, column string) (bool, string) {
	modelFields := reflect.TypeOf(model)

	for i := 0; i < modelFields.NumField(); i++ {
		field := modelFields.Field(i)
		databaseColumn := field.Tag.Get("db")

		if databaseColumn == column {
			return true, databaseColumn
		}
	}

	return false, ""
}
