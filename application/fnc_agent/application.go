package fnc_agent

import (
	"context"
	"github.com/uanid/fakenews-server/application/configs"
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

func NewApplication(cfg *configs.FncConfig) (*App, error) {
	app := &App{}

	err := app.registerServices(cfg)
	if err != nil {
		return nil, err
	}

	app.registerController()

	return app, nil
}

func (a *App) registerServices(cfg *configs.FncConfig) error {
	awsCfg, err := aws_service.NewConfig(cfg.Credentials.ToAwsServiceOption())
	if err != nil {
		return err
	}
	a.ddbService = ddb_service.NewService(*awsCfg, cfg.DynamoDBTable, cfg.DynamoDBRegion)
	sqsService := sqs_service.NewService(*awsCfg, cfg.SqsUrl, cfg.SqsRegion)

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
