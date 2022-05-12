package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/uanid/fakenews-server/pkg/services"
	ddb_service "github.com/uanid/fakenews-server/pkg/services/ddb-service"
	"github.com/uanid/fakenews-server/pkg/types"
)

type Controller struct {
	agentService *services.AgentService
	ddbService   *ddb_service.Service
}

func NewController(pollSvc *services.AgentService, ddbService *ddb_service.Service) *Controller {
	return &Controller{agentService: pollSvc, ddbService: ddbService}
}

func (c *Controller) RunWithTicker(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		err := c.RunOnce(ctx)
		if err != nil {
			return err
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			fmt.Printf("Terminating main controller loop\n")
			return nil
		}
	}
}

func (c *Controller) RunOnce(ctx context.Context) error {
	fmt.Printf("[Agent] Start Polling\n")
	req, ok, err := c.agentService.PollRequest(ctx)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Printf("[Agent] Nothing Polled, Stop RunOnce\n")
		return nil
	}
	fmt.Printf("[Agent] Poll Success uuid='%s'\n", req.Uuid)

	err = c.ddbService.UpdateStatus(ctx, req.Uuid, types.Started)
	if err != nil {
		return fmt.Errorf("UpdateStatusFailed: uuid=%s, %s", req.Uuid, err.Error())
	}

	fmt.Printf("[Agent] Start Run Core Program uuid=%s\n", req.Uuid)
	result, err := c.agentService.RunCore(ctx, req)
	if err == nil {
		err = c.ddbService.UpdateResultWithStatus(ctx, req.Uuid, result)
		if err != nil {
			return fmt.Errorf("UpdateStatusFailed: uuid=%s, %s", req.Uuid, err.Error())
		}
	} else {
		fmt.Printf("[Agent] ErrorOccured while run core: %s\n", err.Error())
		err = c.ddbService.UpdateStatus(ctx, req.Uuid, types.Errored)
		if err != nil {
			return fmt.Errorf("UpdateStatusFailed: uuid=%s, %s", req.Uuid, err.Error())
		}
	}
	fmt.Printf("[Agent] Finished RunOnce uuid=%s\n", req.Uuid)
	return nil
}
