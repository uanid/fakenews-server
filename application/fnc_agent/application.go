package fnc_agent

import (
	"context"
	"time"

	"github.com/uanid/fakenews-server/controllers/agent"
	"github.com/uanid/fakenews-server/pkg/services"
	aws_service "github.com/uanid/fakenews-server/pkg/services/aws-service"
	ddb_service "github.com/uanid/fakenews-server/pkg/services/ddb-service"
	sqs_service "github.com/uanid/fakenews-server/pkg/services/sqs-service"
)

type App struct {
	agent *agent.Controller

	agentService *services.AgentService
	ddbService   *ddb_service.Service
}

func NewApplication(ddbName string, sqsUrl string, profile string, region string) (*App, error) {
	app := &App{}

	err := app.registerServices(ddbName, sqsUrl, profile, region)
	if err != nil {
		return nil, err
	}

	app.registerController()

	return app, nil
}

func (a *App) registerServices(ddbName string, sqsUrl string, profile string, region string) error {
	cfg, err := aws_service.NewConfig(profile, region)
	if err != nil {
		return err
	}
	a.ddbService = ddb_service.NewService(*cfg, ddbName)
	sqsService := sqs_service.NewService(*cfg, sqsUrl)

	a.agentService = services.NewAgentService(a.ddbService, sqsService)
	return nil
}

func (a *App) registerController() {
	a.agent = agent.NewController(a.agentService, a.ddbService)
}

func (a *App) Start(ctx context.Context) error {
	return a.agent.RunOnce(ctx)
}

func (a *App) StartWithTicker(ctx context.Context, interval time.Duration) error {
	return a.agent.RunWithTicker(ctx, interval)
}
