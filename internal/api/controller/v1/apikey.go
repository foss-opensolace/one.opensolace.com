package v1

import (
	"encoding/json"
	"time"

	"github.com/foss-opensolace/api.opensolace.com/internal/api/model/dto"
	"github.com/foss-opensolace/api.opensolace.com/internal/api/service"
	"github.com/foss-opensolace/api.opensolace.com/pkg/middleware"
	"github.com/foss-opensolace/api.opensolace.com/pkg/utils"
	"github.com/foss-opensolace/api.opensolace.com/pkg/validate"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-multierror"
)

func NewAPIKeyRouter(router fiber.Router) {
	group := router.Group("/key")

	group.Post("", middleware.KeyPermission(dto.APIKeyPermissions{KeyCreate: utils.ToPtr(true)}), keyGenerateHandler())
}

func keyGenerateHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body dto.APIKeyCreate

		if err := c.BodyParser(&body); err != nil {
			if e, ok := err.(*json.UnmarshalTypeError); ok {
				return c.Status(fiber.StatusBadRequest).SendString("Cannot use " + e.Value + " from '" + e.Field + "' as a value of type " + e.Type.String())
			}

			if e, ok := err.(*time.ParseError); ok {
				return c.Status(fiber.StatusBadRequest).SendString("Couldn't parse " + e.Value + ". Expected layout: " + e.Layout)
			}

			return err
		}

		if err := validate.Struct(&body); err != nil {
			if errs, ok := err.(*multierror.Error); ok {
				return c.Status(fiber.StatusBadRequest).JSON(errs.Errors)
			}

			return fiber.ErrInternalServerError
		}

		key, err := service.APIKey.Create(body)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(key)
	}
}
