package types

type FakeNewsReq struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	Category int    `json:"category"`
}

type AnalyzeStatus int

const (
	Waiting AnalyzeStatus = iota + 1
	Started
	Finished
)

type AnalyzeResult int

const (
	None AnalyzeResult = iota + 1
	Error
	FakeNews
	TruthNews
)

type AnalyzeRequest struct {
	FakeNews *FakeNewsReq  `dynamodbav:"Fakenews"`
	Uuid     string        `dynamodbav:"Key"`
	Status   AnalyzeStatus `dynamodbav:"Status,int"`
	Result   AnalyzeResult `dynamodbav:"Result,int"`
}
