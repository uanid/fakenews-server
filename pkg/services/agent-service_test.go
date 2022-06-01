package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/uanid/fakenews-server/pkg/types"
)

func TestAgentService_RunCore(t *testing.T) {
	a := &AgentService{}
	s, err := a.RunCore(context.Background(), &types.AnalyzeRequest{
		Uuid: "231dswdas",
		FakeNews: &types.FakeNewsReq{
			Title: "dasdsa",
			Body:  "dwdsd",
		},
		Status: 0,
		Result: "",
	})
	fmt.Printf("Result:%s\n", s)
	fmt.Printf("err: %s\n", err.Error())
}
