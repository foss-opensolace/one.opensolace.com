package v1

import (
	"errors"
	"strconv"

	"github.com/foss-opensolace/api.opensolace.com/internal/api/model/dto"
	"github.com/foss-opensolace/api.opensolace.com/internal/api/service"
	"github.com/foss-opensolace/api.opensolace.com/pkg/middleware"
	"github.com/foss-opensolace/api.opensolace.com/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewUserRouter(router fiber.Router) {
	group := router.Group("/user")

	group.Get(`/id/:id`, middleware.KeyPermission(dto.APIKeyPermissions{UserRead: utils.ToPtr(true)}), userGetOneByIdHandler())
	group.Get("/:username", middleware.KeyPermission(dto.APIKeyPermissions{UserRead: utils.ToPtr(true)}), userGetOneByUsernameHandler())
}

func userGetOneByIdHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Params("id")

		value, err := strconv.ParseUint(userID, 10, 0)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid id provided")
		}

		user, err := service.User.GetById(uint(value))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("User not found with that ID")
			}

			return err
		}

		return c.Status(fiber.StatusOK).JSON(user.ToSafe())
	}
}

func userGetOneByUsernameHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Params("username")

		user, err := service.User.GetByUsername(username)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("User not found with that username")
			}

			return err
		}

		return c.Status(fiber.StatusOK).JSON(user.ToSafe())
	}
}
