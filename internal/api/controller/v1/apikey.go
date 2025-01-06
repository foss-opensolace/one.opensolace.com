package v1

import (
	"github.com/foss-opensolace/api.opensolace.com/internal/api/model/dto"
	"github.com/foss-opensolace/api.opensolace.com/internal/api/service"
	"github.com/foss-opensolace/api.opensolace.com/pkg/middleware"
	"github.com/foss-opensolace/api.opensolace.com/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func NewAPIKeyRouter(router fiber.Router) {
	group := router.Group("/key")

	group.Post("", middleware.KeyPermission(dto.APIKeyPermissions{KeyCreate: utils.ToPtr(true)}), keyGenerateHandler())
}

func keyGenerateHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body dto.APIKeyCreate

		if err := utils.ParseBody(c, &body); err != nil {
			return err
		}

		key, err := service.APIKey.Create(body)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(key)
	}
}
