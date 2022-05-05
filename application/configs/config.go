package configs

import (
	aws_service "github.com/uanid/fakenews-server/pkg/services/aws-service"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type FncConfig struct {
	Credentials    AWSCredential `json:"credentials" yaml:"credentials"`
	DynamoDBTable  string        `json:"dynamoDBTable" yaml:"dynamoDBTable"`
	DynamoDBRegion string        `json:"dynamoDBRegion" yaml:"dynamoDBRegion"`
	SqsUrl         string        `json:"sqsUrl" yaml:"sqsUrl"`
	SqsRegion      string        `json:"sqsRegion" yaml:"sqsRegion"`
}

type AWSCredential struct {
	Profile         string `json:"profile" yaml:"profile"`
	AccessKeyId     string `json:"accessKeyId" yaml:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey" yaml:"secretAccessKey"`
	SessionToken    string `json:"sessionToken" yaml:"sessionToken"`
	StsRegion       string `json:"stsRegion" yaml:"stsRegion"`
}

func LoadConfig(path string) (*FncConfig, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &FncConfig{}
	err = yaml.Unmarshal(buf, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *AWSCredential) ToAwsServiceOption() *aws_service.AwsConfigOption {
	return &aws_service.AwsConfigOption{
		Profile:         c.Profile,
		AccessKeyId:     c.AccessKeyId,
		SecretAccessKey: c.SecretAccessKey,
		SessionToken:    c.SessionToken,
		StsRegion:       c.StsRegion,
	}
}