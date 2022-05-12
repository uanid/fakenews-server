package ddb_service

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/uanid/fakenews-server/pkg/types"
)

type Service struct {
	ddbClient    *dynamodb.Client
	ddbTableName string
}

func NewService(cfg aws.Config, ddbTableName string, region string) *Service {
	copiedCfg := cfg.Copy()
	copiedCfg.Region = region
	return &Service{ddbClient: dynamodb.NewFromConfig(copiedCfg), ddbTableName: ddbTableName}
}

func (s *Service) CreateItem(ctx context.Context, req *types.AnalyzeRequest) error {
	doc, err := attributevalue.MarshalMap(req)
	if err != nil {
		return err
	}

	_, err = s.ddbClient.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      doc,
		TableName: &s.ddbTableName,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetItem(ctx context.Context, uuid string) (*types.AnalyzeRequest, error) {
	out, err := s.ddbClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]ddbTypes.AttributeValue{
			"Key": &ddbTypes.AttributeValueMemberS{Value: uuid},
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

func (s *Service) ListItem(ctx context.Context) ([]types.AnalyzeRequest, error) {
	paginator := dynamodb.NewScanPaginator(s.ddbClient, &dynamodb.ScanInput{
		TableName: &s.ddbTableName,
		Limit:     aws.Int32(1),
	})
	return s.listRequest0(ctx, paginator)
}

func (s *Service) listRequest0(ctx context.Context, paginator *dynamodb.ScanPaginator) ([]types.AnalyzeRequest, error) {
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

func (s *Service) UpdateStatus(ctx context.Context, uuid string, status types.AnalyzeStatus) error {
	_, err := s.ddbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]ddbTypes.AttributeValue{
			"Key": &ddbTypes.AttributeValueMemberS{Value: uuid},
		},
		TableName: &s.ddbTableName,
		AttributeUpdates: map[string]ddbTypes.AttributeValueUpdate{
			"Status": {
				Action: ddbTypes.AttributeActionPut,
				Value: &ddbTypes.AttributeValueMemberN{
					Value: fmt.Sprintf("%d", int(status)),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateResultWithStatus(ctx context.Context, uuid string, result string) error {
	_, err := s.ddbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]ddbTypes.AttributeValue{
			"Key": &ddbTypes.AttributeValueMemberS{Value: uuid},
		},
		TableName: &s.ddbTableName,
		AttributeUpdates: map[string]ddbTypes.AttributeValueUpdate{
			"Status": {
				Action: ddbTypes.AttributeActionPut,
				Value: &ddbTypes.AttributeValueMemberN{
					Value: fmt.Sprintf("%d", int(types.Finished)),
				},
			},
			"Result": {
				Action: ddbTypes.AttributeActionPut,
				Value: &ddbTypes.AttributeValueMemberS{
					Value: result,
				},
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
