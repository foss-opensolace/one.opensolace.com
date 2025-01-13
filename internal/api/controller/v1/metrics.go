package v1

import (
	"github.com/foss-opensolace/one.opensolace.com/internal/api/model/dto"
	"github.com/foss-opensolace/one.opensolace.com/pkg/middleware"
	"github.com/foss-opensolace/one.opensolace.com/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func NewMetricRouter(router fiber.Router) {
	router.Get("", middleware.KeyPermission(dto.APIKeyPermissions{Metrics: utils.True}), handlerGetMetrics())
	router.Get("/health", middleware.KeyPermission(dto.APIKeyPermissions{Health: utils.True}), handlerGetMetricsHealth())
}

func handlerGetMetrics() fiber.Handler {
	return monitor.New(monitor.Config{APIOnly: true})
}

func handlerGetMetricsHealth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Systems healthy!")
	}
}
