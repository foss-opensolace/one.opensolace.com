package v1

import (
	"errors"

	"github.com/foss-opensolace/api.opensolace.com/internal/api/model/dto"
	"github.com/foss-opensolace/api.opensolace.com/internal/api/service"
	"github.com/foss-opensolace/api.opensolace.com/pkg/jwt"
	"github.com/foss-opensolace/api.opensolace.com/pkg/validate"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-multierror"
	"gorm.io/gorm"
)

func NewAuthRouter(router fiber.Router) {
	group := router.Group("/auth")

	group.Post("/register", authRegisterHandler())
	group.Post("/login", authLoginHandler())
}

func authRegisterHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body dto.UserRegister

		if err := c.BodyParser(&body); err != nil {
			return err
		}

		if err := validate.Struct(&body); err != nil {
			if errs, ok := err.(*multierror.Error); ok {
				return c.Status(fiber.StatusBadRequest).JSON(errs.Errors)
			}

			return fiber.ErrInternalServerError
		}

		user, err := service.User.Create(&body)
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
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

func authLoginHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body dto.UserLogin

		if err := c.BodyParser(&body); err != nil {
			return err
		}

		if err := validate.Struct(&body); err != nil {
			if errs, ok := err.(*multierror.Error); ok {
				return c.Status(fiber.StatusBadRequest).JSON(errs.Errors)
			}

			return fiber.ErrInternalServerError
		}

		user, err := service.User.GetByLoginAndPassword(body.Login, body.Password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Username, email or password are incorrect or doesn't exist")
		}

		token, err := jwt.GenerateJWT(user.ID)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"token": token})
	}
}
