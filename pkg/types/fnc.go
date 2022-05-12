package types

type FakeNewsReq struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type AnalyzeStatus int

const (
	Waiting AnalyzeStatus = iota + 1
	Started
	Finished
	Errored
)

type AnalyzeRequest struct {
	Uuid     string        `json:"uuid" dynamodbav:"Key"`
	FakeNews *FakeNewsReq  `json:"fakeNews" dynamodbav:"Fakenews"`
	Status   AnalyzeStatus `json:"status" dynamodbav:"Status,int"`
	Result   string        `json:"result" dynamodbav:"Result"`
}
