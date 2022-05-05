package aws_service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AwsConfigOption struct {
	Profile         string
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
	StsRegion       string
}

func NewConfig(option *AwsConfigOption) (*aws.Config, error) {
	ctx := context.TODO()

	var options []func(*config.LoadOptions) error
	if option.Profile != "" {
		fmt.Printf("[AWS] Use Shared Profile\n")
		options = append(options, config.WithSharedConfigProfile(option.Profile))
	} else {
		fmt.Printf("[AWS] Use Static Credential\n")
		options = append(options, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(option.AccessKeyId, option.SecretAccessKey, option.SessionToken)))
	}
	options = append(options, config.WithRegion(option.StsRegion))

	cfg, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return nil, err
	}
	stsClient := sts.NewFromConfig(cfg)
	stsOut, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}
	fmt.Printf("[AWS] account=%+v, iam-arn=%+v\n", *stsOut.Account, *stsOut.Arn)

	return &cfg, nil
}
