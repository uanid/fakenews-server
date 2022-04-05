package services

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

func (s *RequestService) CreateRequest(ctx context.Context, news *types.FakeNewsReq) (string, error) {
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
