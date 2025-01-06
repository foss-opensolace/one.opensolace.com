package middleware

import (
	"fmt"
	"reflect"

	"github.com/foss-opensolace/api.opensolace.com/internal/api/model/dto"
	"github.com/foss-opensolace/api.opensolace.com/internal/api/service"
	"github.com/foss-opensolace/api.opensolace.com/pkg/exception"
	"github.com/gofiber/fiber/v2"
)

func ValidateKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-API-KEY")

		if key == "" {
			exception.SetID(c, exception.MissingAPIKeyHeader)
			return c.Status(fiber.StatusUnauthorized).SendString("No API key header present (X-API-KEY)")
		}

		apiKey, err := service.APIKey.GetByKey(key)
		if err != nil {
			exception.SetID(c, exception.InvalidAPIKey)
			return c.Status(fiber.StatusUnauthorized).SendString("API key not found")
		}

		if apiKey.CanUse == false {
			exception.SetID(c, exception.CannotUseAPIKey)

			if apiKey.MaxUsage != nil && apiKey.TimesUsed >= *apiKey.MaxUsage {
				return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("API key limit reached: %d/%d", apiKey.TimesUsed, *apiKey.MaxUsage))
			}

			if apiKey.RevokeReason != nil {
				return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("API key has been revoked by the reason: %s", *apiKey.RevokeReason))
			}
		}

		if err := service.APIKey.RegisterUseKey(key); err != nil {
			return err
		}

		c.Locals("apikey", apiKey)
		return c.Next()
	}
}

func KeyPermission(permissions dto.APIKeyPermissions) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Locals("apikey").(*dto.APIKeyLookup)

		permVal := reflect.ValueOf(permissions)
		permType := reflect.TypeOf(permissions)
		keyPermVal := reflect.ValueOf(key.Permissions)

		for i := 0; i < permType.NumField(); i++ {
			field := permType.Field(i)
			permField := permVal.Field(i)
			keyPermField := keyPermVal.FieldByName(field.Name)

			if !permField.IsNil() && !reflect.DeepEqual(*permField.Interface().(*bool), keyPermField.Interface().(bool)) {
				return c.Status(fiber.StatusUnauthorized).SendString("API key doesn't have the proper permissions: " + field.Name)
			}
		}

		return c.Next()
	}
}
