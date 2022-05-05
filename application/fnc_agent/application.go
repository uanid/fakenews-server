package fnc_agent

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/uanid/fakenews-server/controllers/agent"
	"github.com/uanid/fakenews-server/pkg/services"
	ddb_service "github.com/uanid/fakenews-server/pkg/services/ddb-service"
	sqs_service "github.com/uanid/fakenews-server/pkg/services/sqs-service"
)

type App struct {
	agent *agent.Controller

	agentService *services.AgentService
	ddbService   *ddb_service.Service
}

func NewApplication(ddbName string, sqsUrl string) (*App, error) {
	app := &App{}

	err := app.registerServices(ddbName, sqsUrl)
	if err != nil {
		return nil, err
	}

	app.registerController()

	return app, nil
}

func (a *App) registerServices(ddbName string, sqsUrl string) error {
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

	a.ddbService = ddb_service.NewService(cfg, ddbName)
	sqsService := sqs_service.NewService(cfg, sqsUrl)

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
