package services

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	types2 "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
	"github.com/uanid/fakenews-server/pkg/types"
)

type RequestService struct {
	ddbClient    *dynamodb.Client
	ddbTableName string
	sqsClient    *sqs.Client
	sqsQueueUrl  string
}

func NewRequestService(cfg aws.Config, ddbTableName string, sqsQueueUrl string) *RequestService {
	return &RequestService{
		ddbClient:    dynamodb.NewFromConfig(cfg),
		ddbTableName: ddbTableName,
		sqsClient:    sqs.NewFromConfig(cfg),
		sqsQueueUrl:  sqsQueueUrl,
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

	doc, err := attributevalue.MarshalMap(req)
	if err != nil {
		return "", err
	}

	_, err = s.ddbClient.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      doc,
		TableName: &s.ddbTableName,
	})
	if err != nil {
		return "", err
	}

	_, err = s.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:            aws.String(id.String()),
		QueueUrl:               &s.sqsQueueUrl,
		MessageDeduplicationId: aws.String(id.String()),
		MessageGroupId:         aws.String("fakenews-analyze"),
	})
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (s *RequestService) GetRequest(ctx context.Context, requestId string) (*types.AnalyzeRequest, error) {
	out, err := s.ddbClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types2.AttributeValue{
			"Key": &types2.AttributeValueMemberS{Value: requestId},
		},
		TableName: &s.ddbTableName,
	})
	if err != nil {
		return nil, err
	}

	req := &types.AnalyzeRequest{}
	err = attributevalue.UnmarshalMap(out.Item, req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *RequestService) ListRequests(ctx context.Context) ([]types.AnalyzeRequest, error) {
	paginator := dynamodb.NewScanPaginator(s.ddbClient, &dynamodb.ScanInput{
		TableName: &s.ddbTableName,
		Limit:     aws.Int32(1),
	})
	return s.listRequest0(ctx, paginator)
}

func (s *RequestService) listRequest0(ctx context.Context, paginator *dynamodb.ScanPaginator) ([]types.AnalyzeRequest, error) {
	out, err := paginator.NextPage(ctx)
	if err != nil {
		return nil, err
	}

	requests := []types.AnalyzeRequest{}
	err = attributevalue.UnmarshalListOfMaps(out.Items, &requests)
	if err != nil {
		return nil, err
	}

	if paginator.HasMorePages() {
		nextRequests, err := s.listRequest0(ctx, paginator)
		if err != nil {
			return nil, err
		}
		return append(requests, nextRequests...), nil
	} else {
		return requests, nil
	}

}
