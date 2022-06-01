package execute_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"syscall"

	"github.com/uanid/fakenews-server/pkg/types"
)

type Input struct {
	Title string `json:"Title"`
	Body  string `json:"Body"`
}

type Output struct {
	TF      any   `json:"TF"`
	Score   int   `json:"Score"`
	Words   []any `json:"Words"`
	Weights []any `json:"Weights"`
}

func WriteInputFile(req *types.FakeNewsReq) (string, error) {
	input := &Input{
		Title: req.Title,
		Body:  req.Body,
	}

	buf, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("input-%04d.json", rand.Uint32()%1000)
	err = ioutil.WriteFile(fileName, buf, 0644)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func GenerateOutputFileName() string {
	return fmt.Sprintf("output-%04d.json", rand.Uint32()%1000)
}

func ReadOutputFile(fileName string) (*Output, error) {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("ReadOutputFileError, fileName:%s, %s", fileName, err.Error())
	}

	out := &Output{}
	err = json.Unmarshal(buf, out)
	if err != nil {
		return nil, fmt.Errorf("ReadOutputFileError, fileName:%s, %s", fileName, err.Error())
	}

	return out, nil
}

func ReadOutputFile2(fileName string) (string, error) {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("ReadOutputFileError, fileName:%s, %s", fileName, err.Error())
	}

	return string(buf), nil
}

func CleanupFiles(fileNames ...string) {
	for _, fileName := range fileNames {
		_ = os.Remove(fileName)
	}
}

func Execute(ctx context.Context, bin string, args ...string) (stdout string, stderr string, exit int, err error) {
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

	wd, err := os.Getwd()
	fmt.Printf("Command execute: %s %v, workdir: %s, exitcode: %d\n", bin, args, wd, exit)
	return
}
