package repositories

type AnalyzeRepository interface {
	InsertAnalyze(title string, body string) (id string, err error)
	GetAnalyze(id string) (status AnalyzeStatus, result string, err error)

	AcquireAnalyze() (id string, title string, body string)
	FinishAnalyze(id string, result string) error
}

type AnalyzeStatus int

const (
	Waiting AnalyzeStatus = iota + 1
	Started
	Finished
)
