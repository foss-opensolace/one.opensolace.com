package validate

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ruleUsername(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[a-z0-9]{3,}$`).Match([]byte(fl.Field().String()))
}

func rulePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 6 {
		return false
	}

	if len([]byte(password)) > 72 {
		return false
	}

	return true
}
