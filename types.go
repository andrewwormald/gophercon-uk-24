package example

type Gopher struct {
	Name   string
	Colour string
	Age    int
	School string
	Work   string
}

//go:generate stringer -type=Status
type Status int

const (
	StepUnknown          Status = 0
	StatusStarted        Status = 1
	StatusNameCreated    Status = 2
	StatusColourSet      Status = 3
	StatusAgeDefined     Status = 4
	StatusSentToSchool   Status = 5
	StatusSentToWork     Status = 6
	StatusFinishedSchool Status = 7
)

type GraduationResponse struct {
	Graduated bool
}
