package application

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/uanid/fakenews-server/controllers"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	fiberApp *fiber.App
	port int
}

func NewApplication(port int) (*App,error) {
	app := &App{}
	app.fiberApp = fiber.New()
	app.fiberApp.Use(requestid.New(), logger.New(), cors.New())

	app.registerControllers()

	app.port = port
	return app
}

func (a *App) registerControllers() {
	a.fiberApp.Get("/api/v1/ping", controllers.Ping)
	a.fiberApp.Post("/api/v1/fakenews-analyze", controllers.RequestAnalyze)
	a.fiberApp.Get("/api/v1/fakenews-analyze/:id", controllers.GetAnalyze)
}

func (a *App) Start() error {
	return a.fiberApp.Listen(fmt.Sprintf(":%d", a.port))
}