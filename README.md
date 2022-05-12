# Fakenews-Server

An API Server for Run The FNC-1 Model

# Start

You should run both `api server` and `agent`.

---------------------------

## Run requirements

You should select one of options.

#### Option: Use AWS CLI Credentials

1. Install AWS Cli [Download Here](https://docs.aws.amazon.com/ko_kr/cli/latest/userguide/getting-started-install.html)
2. Setup IAM Credential `$ aws configure --profile=fnc`

#### Option: Use Plain String Credentials

Write aws credentials like below

```yaml
credentials:
  accessKeyId: DWDSESDE5N7I
  secretAccessKey: dw32depoH+RmohMHhTPEc7HeNe0zRuvA2TrreK
  stsRegion: ap-northeast-2
```

---------------------------

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

## Global Parameter

```shell
# --config=<configPath> : Application config file
```

---------------------------

## Server API

| Method | Url                          | Body                          |
|--------|------------------------------|-------------------------------|
| GET    | /api/v1/ping                 |                               |
| POST   | /api/v1/fakenews-analyze/    | {"title":"aaa","body":"bbbb"} |
| GET    | /api/v1/fakenews-analyze/:id |                               |
| GET    | /api/v1/fakenews-analyze/    |                               |

## Server API Example
```shell
curl http://localhost:8080/api/v1/ping

curl http://localhost:8080/api/v1/fakenews-analyze/

curl http://localhost:8080/api/v1/fakenews-analyze/71a09766-d1a8-11ec-b4f4-c26f8fc1fbdc

curl -X POST -H "Content-Type: application/json" http://localhost:8080/api/v1/fakenews-analyze -d @curlinput.json
```
