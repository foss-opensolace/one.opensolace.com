package v1

import (
	"errors"

	"github.com/foss-opensolace/one.opensolace.com/internal/api/model/dto"
	"github.com/foss-opensolace/one.opensolace.com/internal/api/service"
	"github.com/foss-opensolace/one.opensolace.com/pkg/exception"
	"github.com/foss-opensolace/one.opensolace.com/pkg/jwt"
	"github.com/foss-opensolace/one.opensolace.com/pkg/middleware"
	"github.com/foss-opensolace/one.opensolace.com/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewAuthRouter(router fiber.Router) {
	router.Post("/register", middleware.KeyPermission(dto.APIKeyPermissions{UserAuthRegister: utils.True}), handlerPostAuthRegister())
	router.Post("/login", middleware.KeyPermission(dto.APIKeyPermissions{UserAuthLogin: utils.True}), handlerPostAuthLogin())
}

func handlerPostAuthRegister() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body dto.UserRegister

		if err := utils.ParseBody(c, &body); err != nil {
			return err
		}

		user, err := service.User.Create(&body)
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				exception.SetID(c, exception.IdAuthDuplicated)
				return c.Status(fiber.StatusConflict).SendString("A user with that username or email already exists")
			}

			return err
		}

		token, err := jwt.GenerateJWT(user.ID)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": token})
	}
}

func handlerPostAuthLogin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body dto.UserLogin

		if err := utils.ParseBody(c, &body); err != nil {
			return err
		}

		user, err := service.User.GetByLoginAndPassword(body.Login, body.Password)
		if err != nil {
			exception.SetID(c, exception.IdAuthInvalidCredentials)
			return c.Status(fiber.StatusUnauthorized).SendString("Username, email or password are incorrect or doesn't exist")
		}

		token, err := jwt.GenerateJWT(user.ID)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"token": token})
	}
}
