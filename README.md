# Fakenews-Server
An API Server for Run The FNC-1 Model

# Start

You should run both `api server` and `agent`.

## Run requirements
1. Install AWS Cli [Download Here](https://docs.aws.amazon.com/ko_kr/cli/latest/userguide/getting-started-install.html)
2. Setup IAM Credential `$ aws configure --profile=fnc`

## Run API Server
```shell
$ fakenews-server server

# Parameters
#--port=<int> : Rest API Listen Port, (Default 8080)
```

## Run Agent
```shell
$ fakenews-server agent

# Parameters
#--interval=<duration> : The Agent Work Iteration Loop Interval (Default 10s)
#--once : Flag for not Loop Iteration (Default false=not-set)
```

## Global Pramaeters
```
--ddb=<table-name> : Dynamodb Table Name (Default fnc1-db)
--sqs=<url> : Full SQS Queue Url
--profile=<aws-cli-profile> : AWS Shared Credential Profile Name (Default fnc)
--region=<aws-region> : Where DynamoDB and SQS located AWS Region Name (Default ap-northeast-2)
```