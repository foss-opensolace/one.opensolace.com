package v1

import (
	"github.com/foss-opensolace/api.opensolace.com/internal/api/model/dto"
	"github.com/foss-opensolace/api.opensolace.com/pkg/middleware"
	"github.com/foss-opensolace/api.opensolace.com/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func NewMetricRouter(router fiber.Router) {
	group := router.Group("/metrics")

	group.Get("", middleware.KeyPermission(dto.APIKeyPermissions{Metrics: utils.ToPtr(true)}), metrics())
	group.Get("/health", middleware.KeyPermission(dto.APIKeyPermissions{Health: utils.ToPtr(true)}), health())
}

func metrics() fiber.Handler {
	return monitor.New(monitor.Config{APIOnly: true})
}

func health() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Systems healthy!")
	}
}
