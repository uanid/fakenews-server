package services

import (
	"context"
	"fmt"
	"strconv"

	ddb_service "github.com/uanid/fakenews-server/pkg/services/ddb-service"
	execute_service "github.com/uanid/fakenews-server/pkg/services/execute-service"
	sqs_service "github.com/uanid/fakenews-server/pkg/services/sqs-service"
	"github.com/uanid/fakenews-server/pkg/types"
)

type AgentService struct {
	ddbService *ddb_service.Service
	sqsService *sqs_service.Service
}

func NewAgentService(ddbService *ddb_service.Service, sqsService *sqs_service.Service) *AgentService {
	return &AgentService{
		ddbService: ddbService,
		sqsService: sqsService,
	}
}

func (s *AgentService) PollRequest(ctx context.Context) (*types.AnalyzeRequest, bool, error) {
	uuid, ok, err := s.sqsService.PollUuid(ctx)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}

	req, err := s.ddbService.GetItem(ctx, uuid)
	if err != nil {
		return nil, false, err
	}
	return req, true, nil
}

func (s *AgentService) RunCore(ctx context.Context, req *types.AnalyzeRequest) (string, error) {
	fmt.Printf("[Core] Execute Command...\n")

	inputFile, err := execute_service.WriteInputFile(req.FakeNews)
	if err != nil {
		return "", err
	}
	outputFile := execute_service.GenerateOutputFileName()
	wordsLen := 5
	defer execute_service.CleanupFiles(inputFile, outputFile)

	stdout, stderr, exitCode, err := execute_service.Execute(ctx, "python3", "FNdwithjson.py", inputFile, outputFile, strconv.Itoa(wordsLen))
	if stdout != "" {
		fmt.Println("[Stdout] " + stdout)
	}
	if stderr != "" {
		fmt.Println("[Stderr] " + stderr)
	}

	if err != nil {
		return "", err
	}
	if exitCode != 0 {
		return "", fmt.Errorf("UnexpectedExitcode: exitcode is not 0, %d", exitCode)
	}

	result, err := execute_service.ReadOutputFile2(outputFile)
	if err != nil {
		return "", err
	}

	fmt.Printf("[Core] Execute Finished result=%s\n", result)
	return result, nil
}
