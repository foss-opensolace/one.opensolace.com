package validate

import "github.com/go-playground/validator/v10"

var instance *validator.Validate

func New() {
	instance = register(validator.New())
}

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

	return v
}
