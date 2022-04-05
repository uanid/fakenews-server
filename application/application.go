package application

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/uanid/fakenews-server/controllers"
	"github.com/uanid/fakenews-server/pkg/services"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	// Controller Layer
	fiberApp *fiber.App
	port     int

	// Service Layer
	requestSvc *services.RequestService
}

func NewApplication(port int) (*App, error) {
	app := &App{}
	app.fiberApp = fiber.New()
	app.fiberApp.Use(requestid.New(), logger.New(), cors.New())

	err := app.registerServices()
	if err != nil {
		return nil, err
	}
	app.registerControllers()

	app.port = port
	return app, nil
}

func (a *App) registerServices() error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("fnc"), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return err
	}

	stsClient := sts.NewFromConfig(cfg)
	stsOut, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("account=%+v, iamarn=%+v\n", *stsOut.Account, *stsOut.Arn)

	a.requestSvc = services.NewRequestService(cfg, "fnc1-db", "https://sqs.ap-northeast-2.amazonaws.com/031804216199/fnc1-queue.fifo")
	return nil
}

func (a *App) registerControllers() {
	controllers.SetRequestSvc(a.requestSvc)

	a.fiberApp.Get("/api/v1/ping", controllers.Ping)
	a.fiberApp.Post("/api/v1/fakenews-analyze", controllers.RequestAnalyze)
	a.fiberApp.Get("/api/v1/fakenews-analyze/:id", controllers.GetAnalyze)
}

func (a *App) Start() error {
	return a.fiberApp.Listen(fmt.Sprintf(":%d", a.port))
}
