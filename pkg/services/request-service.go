package services

import (
	"context"

	"github.com/google/uuid"
	ddb_service "github.com/uanid/fakenews-server/pkg/services/ddb-service"
	sqs_service "github.com/uanid/fakenews-server/pkg/services/sqs-service"
	"github.com/uanid/fakenews-server/pkg/types"
)

type RequestService struct {
	ddbService *ddb_service.Service
	sqsService *sqs_service.Service
}

func NewRequestService(ddbService *ddb_service.Service, sqsService *sqs_service.Service) *RequestService {
	return &RequestService{
		ddbService: ddbService,
		sqsService: sqsService,
	}
}

func (s *RequestService) CreateAnalyze(ctx context.Context, news *types.FakeNewsReq) (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	req := &types.AnalyzeRequest{
		FakeNews: news,
		Uuid:     id.String(),
		Status:   types.Waiting,
	}

	err = s.ddbService.CreateItem(ctx, req)
	if err != nil {
		return "", err
	}

	err = s.sqsService.SendUuid(ctx, id.String())
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (s *RequestService) GetRequest(ctx context.Context, requestId string) (*types.AnalyzeRequest, error) {
	return s.ddbService.GetItem(ctx, requestId)
}

func (s *RequestService) ListRequests(ctx context.Context) ([]types.AnalyzeRequest, error) {
	return s.ddbService.ListItem(ctx)
}
