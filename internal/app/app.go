package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/foss-opensolace/api.opensolace.com/internal/api/controller"
	"github.com/foss-opensolace/api.opensolace.com/internal/config"
	"github.com/foss-opensolace/api.opensolace.com/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type server struct {
	instance *fiber.App
}

func New() {
	app := server{instance: fiber.New(fiber.Config{Immutable: true})}
	app.loadMiddlewares()
	controller.New(app.instance)
	app.listen()
}

func (s *server) loadMiddlewares() {
	s.instance.Use(
		middleware.RequestId(),
		middleware.Recover(),
		middleware.CORS(),
		middleware.Helmet(),
		middleware.Logger(),
		middleware.Interceptor(),
		middleware.ValidateKey(),
	)
}

func (s *server) listen() {
	go func() {
		port, exists := os.LookupEnv("PORT")

		if !exists {
			port = config.Router.Port
		}

		if err := s.instance.Listen(fmt.Sprintf("%s:%s", config.Router.Host, port)); err != nil {
			panic(err)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
}
