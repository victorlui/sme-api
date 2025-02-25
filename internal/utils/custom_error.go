package utils

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

func CustomErrorHandling(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			fmt.Println(fieldErr.Namespace())
			fmt.Println(fieldErr.Field())
			switch fieldErr.Tag() {
			case "required":
				errors[fieldErr.Field()] = fieldErr.Field() + " é obrigatório."
			case "min":
				errors[fieldErr.Field()] = fieldErr.Field() + " deve ter pelo menos " + fieldErr.Param() + " caracteres."
			case "email":
				errors[fieldErr.Field()] = "O campo " + fieldErr.Field() + " deve ser um e-mail válido."
			case "ISO8601date":
				errors[fieldErr.Field()] = "O campo " + fieldErr.Field() + " deve ser uma data válida."
			default:
				errors[fieldErr.Field()] = "Erro no campo " + fieldErr.Field() + "."
			}
		}
	}
	return errors
}

func CustomDate(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().Interface().(time.Time).Format("2006-01-02"))
	fmt.Println("aqui", err)
	return err == nil
}

func IsISO8601Date(fl validator.FieldLevel) bool {
	ISO8601DateRegexString := "^(?:[1-9]\\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d(?:\\.\\d{1,9})?(?:Z|[+-][01]\\d:[0-5]\\d)$"
	ISO8601DateRegex := regexp.MustCompile(ISO8601DateRegexString)
	return ISO8601DateRegex.MatchString(fl.Field().String())
}
