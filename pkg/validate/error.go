package validate

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-multierror"
)

type ParamError struct {
	Param   string `json:"param"`
	Message string `json:"message"`
}

func (pe ParamError) Error() string {
	return fmt.Sprintf("%s: %s", pe.Param, pe.Message)
}

func toParamError(fe validator.FieldError) ParamError {
	field := fe.Field()
	var msg string

	switch fe.Tag() {
	case "required":
		msg = "Is empty, but is required"
	case "eqfield":
		msg = fmt.Sprintf("Doesn't match %s", fe.Param())
	case "min":
		msg = fmt.Sprintf("%s characters minimum", fe.Param())
	case "max":
		msg = fmt.Sprintf("%s characters maximum", fe.Param())
	case "username":
		msg = "Usernames must be composed of lowercase alphanumerical characters"
	case "email":
		msg = "Incorrect email format"
	case "password":
		msg = "Password must contain at least 6 digits and cannot exceed 72 bytes. Reduce the length or avoid characters like emojis."
	default:
		msg = fe.Error()
	}

	return ParamError{
		Param:   field,
		Message: msg,
	}
}

func getStructErrors(s any) validator.ValidationErrors {
	if err := instance.Struct(s); err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

func multiError(errs []ParamError) error {
	var result *multierror.Error

	for _, e := range errs {
		result = multierror.Append(result, e)
	}

	return result.ErrorOrNil()
}
