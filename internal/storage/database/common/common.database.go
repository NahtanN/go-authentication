package database

import (
	"reflect"
	"time"

	"github.com/nahtann/go-authentication/internal/utils"
)

type QueryData struct {
	SearchFields []string
	SearchArgs   []any
}

func SetQueryData(model interface{}) (*QueryData, error) {
	modelTypes := reflect.TypeOf(model)
	modelValues := reflect.ValueOf(model)

	if modelTypes.Kind() != reflect.Struct {
		return nil, &utils.CustomError{
			Message: "Invalid struct model.",
		}
	}

	qd := QueryData{}

	for i := 0; i < modelTypes.NumField(); i++ {
		modelField := modelTypes.Field(i)

		field := modelField.Tag.Get("db")
		value := modelValues.Field(i).Interface()

		if modelField.Type.Kind() == reflect.Uint32 && value.(uint32) == 0 {
			continue
		}

		if modelField.Type == reflect.TypeOf(time.Time{}) && value.(time.Time).IsZero() {
			continue
		}

		if field != "" && value != "" {
			qd.SearchFields = append(qd.SearchFields, field)
			qd.SearchArgs = append(qd.SearchArgs, value)
		}
	}

	return &qd, nil
}
