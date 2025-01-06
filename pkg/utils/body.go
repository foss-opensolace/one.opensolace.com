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
			exception.SetID(c, exception.InvalidFieldType)
			return c.Status(fiber.StatusBadRequest).SendString("Cannot use " + e.Value + " from '" + e.Field + "' as a value of type " + e.Type.String())
		}

		if e, ok := err.(*time.ParseError); ok {
			exception.SetID(c, exception.InvalidFieldLayout)
			return c.Status(fiber.StatusBadRequest).SendString("Couldn't parse " + e.Value + ". Expected layout: " + e.Layout)
		}

		return err
	}

	if err := validate.Struct(&out); err != nil {
		if errs, ok := err.(*multierror.Error); ok {
			return c.Status(fiber.StatusBadRequest).JSON(errs.Errors)
		}

		return fiber.ErrInternalServerError
	}

	return nil
}
