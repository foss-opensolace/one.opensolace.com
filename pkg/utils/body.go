package utils

import (
	"encoding/json"
	"time"

	"github.com/foss-opensolace/api.opensolace.com/pkg/exception"
	"github.com/foss-opensolace/api.opensolace.com/pkg/validate"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-multierror"
)

func ParseBody(c *fiber.Ctx, out any) error {
	if err := c.BodyParser(&out); err != nil {
		if e, ok := err.(*json.UnmarshalTypeError); ok {
			return exception.FieldTypeError{Value: e.Value, Field: e.Field, Type: e.Type.String()}
		}

		if e, ok := err.(*time.ParseError); ok {
			return exception.FieldLayoutError{Value: e.Value, Layout: e.Layout}
		}

		return err
	}

	if err := validate.Struct(out); err != nil {
		if errs, ok := err.(*multierror.Error); ok {
			return errs
		}

		return err
	}

	return nil
}
