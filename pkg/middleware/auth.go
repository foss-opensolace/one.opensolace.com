package middleware

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/foss-opensolace/one.opensolace.com/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

func validateUserToken(c *fiber.Ctx) error {
	tokenString := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("No JWT token found")
	}

	claims, err := jwt.GetClaimsJWT(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	id, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil || id < 0 {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid id from JWT")
	}

	c.Locals("user_id", fmt.Sprint(id))

	return nil
}

func Authorize() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := validateUserToken(c); err != nil {
			return err
		}

		return c.Next()
	}
}

func OptionalAuthorize() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := validateUserToken(c); err != nil {
			return c.Next()
		}

		return c.Next()
	}
}
