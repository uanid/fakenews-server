package sqs_service

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const messageGroupId = "fakenews-analyze"

type Service struct {
	sqsClient   *sqs.Client
	sqsQueueUrl string
}

func NewService(cfg aws.Config, sqsQueueUrl string) *Service {
	return &Service{
		sqsClient:   sqs.NewFromConfig(cfg),
		sqsQueueUrl: sqsQueueUrl,
	}
}

func (s *Service) PollUuid(ctx context.Context) (string, bool, error) {
	out, err := s.sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.sqsQueueUrl),
		AttributeNames:      []types.QueueAttributeName{types.QueueAttributeNameAll},
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     5,
		//VisibilityTimeout:       100,
	})
	if err != nil {
		return "", false, err
	}

	if len(out.Messages) > 1 {
		return "", false, fmt.Errorf("UnexpectedPollError: returned message count should be less than 1, count=%d", len(out.Messages))
	}

	if len(out.Messages) == 0 {
		return "", false, nil
	}

	message := out.Messages[0]

	_, err = s.sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.sqsQueueUrl),
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		return "", false, fmt.Errorf("UnexpectedPollError: delete message failed, %s", err.Error())
	}
	return *message.Body, true, err
}

func (s *Service) SendUuid(ctx context.Context, uuid string) error {
	_, err := s.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:            aws.String(uuid),
		QueueUrl:               &s.sqsQueueUrl,
		MessageDeduplicationId: aws.String(uuid),
		MessageGroupId:         aws.String(messageGroupId),
	})
	if err != nil {
		return err
	}
	return nil
}
