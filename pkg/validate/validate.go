package validate

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var instance *validator.Validate = register(validator.New())

func Struct(s any) error {
	var errors []ParamError

	for _, v := range getStructErrors(s) {
		errors = append(errors, toParamError(v))
	}

	return multiError(errors)
}

func register(v *validator.Validate) *validator.Validate {
	v.RegisterValidation("username", ruleUsername)
	v.RegisterValidation("password", rulePassword)

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return v
}
