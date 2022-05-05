package services

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	ddb_service "github.com/uanid/fakenews-server/pkg/services/ddb-service"
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

func (s *AgentService) RunCore(ctx context.Context, req *types.AnalyzeRequest) (types.AnalyzeResult, error) {
	fmt.Printf("[Core] Execute Command...\n")

	stdout, stderr, exitCode, err := execute(ctx, "python3", "--help")
	if err != nil {
		return types.Error, err
	}
	if exitCode != 0 {
		return types.Error, fmt.Errorf("UnexpectedExitcode: exitcode is not 0, %d", exitCode)
	}

	if stdout != "" {
		fmt.Println("[Stdout] " + stdout)
	}

	if stderr != "" {
		fmt.Println("[Stderr] " + stderr)
	}

	var result types.AnalyzeResult
	if strings.Contains(stdout, "") {
		result = types.TruthNews
	} else {
		result = types.FakeNews
	}

	fmt.Printf("[Core] Execute Finished result=%d\n", int(result))
	return result, nil
}

func execute(ctx context.Context, bin string, args ...string) (stdout string, stderr string, exit int, err error) {
	exit = 0

	cmd := exec.CommandContext(ctx, bin, args...)
	var bufOut bytes.Buffer
	var bufErr bytes.Buffer
	cmd.Stdout = &bufOut
	cmd.Stderr = &bufErr

	cmdErr := cmd.Run()
	if cmdErr != nil {
		if exiterr, ok := cmdErr.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exit = status.ExitStatus()
			} else {
				err = fmt.Errorf("ExitError: %s", cmdErr.Error())
			}
		} else if exiterr2, ok := cmdErr.(*exec.Error); ok {
			err = fmt.Errorf("ExecError: PATH=%s, %s", os.Getenv("PATH"), exiterr2.Error())
		} else {
			err = fmt.Errorf("ExecUnknownError: %s", cmdErr.Error())
		}
	}

	stdout = bufOut.String()
	stderr = bufErr.String()
	fmt.Printf("Command execute: %s %v, exitcode: %d\n", bin, args, exit)
	return
}
