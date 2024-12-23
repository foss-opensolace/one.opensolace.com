package validate

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-multierror"
)

type ParamError struct {
	Param   string
	Message string
}

func (pe ParamError) Error() string {
	return fmt.Sprintf("%s: %s", pe.Param, pe.Message)
}

func toParamError(fe validator.FieldError) ParamError {
	field := strings.ToLower(fe.Field())

	switch fe.Tag() {
	case "username":
		return ParamError{
			Param:   field,
			Message: "Usernames must be composed of lowercase alphanumerical characters",
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
		if e, ok := err.(*validator.InvalidValidationError); ok {
			panic(e)
		}

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
