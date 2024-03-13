package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Validate(s any) string {
	validate := validator.New()

	err := validate.Struct(s)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)

		errorMessages := []string{}
		for _, e := range validationErrors {
			structMessage := fmt.Sprintf(
				"Field `%s` failed validation.",
				e.Field(),
			)

			errorMessages = append(errorMessages, structMessage)
		}

		return strings.Join(errorMessages, " ")
	}

	return ""
}
