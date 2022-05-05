package fnc_server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/uanid/fakenews-server/controllers/rest"
	"github.com/uanid/fakenews-server/pkg/services"
	aws_service "github.com/uanid/fakenews-server/pkg/services/aws-service"
	ddb_service "github.com/uanid/fakenews-server/pkg/services/ddb-service"
	sqs_service "github.com/uanid/fakenews-server/pkg/services/sqs-service"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	// Controller Layer
	fiberApp *fiber.App
	port     int

	// Service Layer
	requestSvc *services.RequestService
}

func NewApplication(port int, ddbName string, sqsUrl string, profile string, region string) (*App, error) {
	app := &App{}
	app.fiberApp = fiber.New()
	app.fiberApp.Use(requestid.New(), logger.New(), cors.New())

	err := app.registerServices(ddbName, sqsUrl, profile, region)
	if err != nil {
		return nil, err
	}
	app.registerControllers()

	app.port = port
	return app, nil
}

func (a *App) registerServices(ddbName string, sqsUrl string, profile string, region string) error {
	cfg, err := aws_service.NewConfig(profile, region)
	if err != nil {
		return err
	}

	ddbService := ddb_service.NewService(*cfg, ddbName)
	sqsService := sqs_service.NewService(*cfg, sqsUrl)

	a.requestSvc = services.NewRequestService(ddbService, sqsService)
	return nil
}

func (a *App) registerControllers() {
	restCtrl := rest.NewRestController(a.requestSvc)

	a.fiberApp.Get("/api/v1/ping", restCtrl.Ping)
	a.fiberApp.Post("/api/v1/fakenews-analyze", restCtrl.RequestAnalyze)
	a.fiberApp.Get("/api/v1/fakenews-analyze/:id", restCtrl.GetAnalyze)
	a.fiberApp.Get("/api/v1/fakenews-analyze", restCtrl.ListAnalyze)
}

func (a *App) Start() error {
	return a.fiberApp.Listen(fmt.Sprintf(":%d", a.port))
}
