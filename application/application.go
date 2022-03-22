package application

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/uanid/fakenews-server/controllers"
)

type App struct {
	fiberApp *fiber.App
	port int
}

func NewApplication(port int) *App {
	app := &App{}
	app.fiberApp = fiber.New()
	app.fiberApp.Use(logger.New())

	app.registerControllers()

	app.port = port
	return app
}

func (a *App) registerControllers() {
	a.fiberApp.Get("/api/v1/ping", controllers.Ping)
}

func (a *App) Start() error {
	return a.fiberApp.Listen(fmt.Sprintf(":%d", a.port))
}