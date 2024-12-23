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

	switch fe.Tag() {
	case "required":
		return ParamError{
			Param:   field,
			Message: "Is empty, but is required",
		}
	case "eqfield":
		return ParamError{
			Param:   field,
			Message: fmt.Sprintf("Doesn't match %s", fe.Param()),
		}
	case "min":
		return ParamError{
			Param:   field,
			Message: fmt.Sprintf("%s characters minimum", fe.Param()),
		}
	case "max":
		return ParamError{
			Param:   field,
			Message: fmt.Sprintf("%s characters maximum", fe.Param()),
		}
	case "username":
		return ParamError{
			Param:   field,
			Message: "Usernames must be composed of lowercase alphanumerical characters",
		}
	case "email":
		return ParamError{
			Param:   field,
			Message: "Incorrect email format",
		}
	case "password":
		return ParamError{
			Param:   field,
			Message: "Password must contain at least 6 digits and cannot exceed 72 bytes. Reduce the length or avoid characters like emojis.",
		}
	default:
		return ParamError{
			Param:   field,
			Message: fe.Error(),
		}
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
