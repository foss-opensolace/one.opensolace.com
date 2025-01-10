package controller

import (
	v1 "github.com/foss-opensolace/api.opensolace.com/internal/api/controller/v1"
	"github.com/gofiber/fiber/v2"
)

type controller struct {
	instance *fiber.App
}

func New(app *fiber.App) {
	c := controller{
		instance: app,
	}

	c.v1()
}

func (c *controller) v1() {
	router := c.instance.Group("/v1")

	v1.NewMetricRouter(router.Group("/metrics"))
	v1.NewAuthRouter(router.Group("/auth"))
	v1.NewUserRouter(router.Group("/user"))
}
