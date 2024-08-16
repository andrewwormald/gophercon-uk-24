package example

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"k8s.io/utils/clock"

	"github.com/luno/workflow"
	"github.com/luno/workflow/adapters/memrecordstore"
	"github.com/luno/workflow/adapters/memrolescheduler"
	"github.com/luno/workflow/adapters/memstreamer"
	"github.com/luno/workflow/adapters/memtimeoutstore"
)

//go:generate go run ./cmd/visualiser

func NewGopherWorkflow(age int, clock clock.Clock) *workflow.Workflow[Gopher, Status] {
	b := workflow.NewBuilder[Gopher, Status]("GopherCon UK - 2024")
	b.AddStep(
		StatusStarted,
		func(ctx context.Context, r *workflow.Run[Gopher, Status]) (Status, error) {
			// Set the name of the gopher
			r.Object.Name = "Budgie"

			// Move to StatusNameCreated
			return StatusNameCreated, nil
		},
		StatusNameCreated,
	)

	b.AddStep(
		StatusNameCreated,
		func(ctx context.Context, r *workflow.Run[Gopher, Status]) (Status, error) {
			// Set the gopher's colour to blue
			r.Object.Colour = "Blue"

			// Move to StatusColourSet
			return StatusColourSet, nil
		},
		StatusColourSet,
	)

	b.AddStep(
		StatusColourSet,
		func(ctx context.Context, r *workflow.Run[Gopher, Status]) (Status, error) {
			// Set the gopher's age
			r.Object.Age = age

			// Move to StatusAgeDefined
			return StatusAgeDefined, nil
		},
		StatusAgeDefined,
	)

	b.AddStep(
		StatusAgeDefined,
		func(ctx context.Context, r *workflow.Run[Gopher, Status]) (Status, error) {
			var destination Status
			if r.Object.Age > 18 {
				r.Object.Work = "Finding lost goroutines"
				destination = StatusSentToWork
			} else {
				r.Object.School = "Learning about the mysterious runtime scheduler"
				destination = StatusSentToSchool
			}
			return destination, nil
		},
		StatusSentToSchool, StatusSentToWork,
	)

	b.AddTimeout(
		StatusSentToSchool,
		func(ctx context.Context, r *workflow.Run[Gopher, Status], now time.Time) (time.Time, error) {
			yearsTillEighteen := 18 - r.Object.Age

			// Wait until the Gopher turns 18 and then execute the below consumer function.
			return now.AddDate(yearsTillEighteen, 0, 0), nil
		},
		func(ctx context.Context, r *workflow.Run[Gopher, Status], now time.Time) (Status, error) {
			r.Object.School = "Graduated school and ready to find lost goroutines"

			// Move to StatusSentToWork now that the Gopher has finished school
			return StatusFinishedSchool, nil
		},
		StatusFinishedSchool,
	)

	b.AddCallback(
		StatusFinishedSchool,
		func(ctx context.Context, r *workflow.Run[Gopher, Status], reader io.Reader) (Status, error) {
			b, err := io.ReadAll(reader)
			if err != nil {
				return 0, err
			}

			var gr GraduationResponse
			err = json.Unmarshal(b, &gr)
			if err != nil {
				return 0, err
			}

			if !gr.Graduated {
				return r.Skip()
			}

			r.Object.Work = "Sent to start finding lost goroutines"

			// Move to StatusSentToWork now that the Gopher has finished school
			return StatusSentToWork, nil
		},
		StatusSentToWork,
	)

	// Using in-memory implementations of the required adapters. Usually this can be passed into the function
	// that builds the workflow to allow for tests to use in-memory versions while staging, and production
	// can use the actual implementations. If you have access to your adapters dependencies in tests then
	// that is advised to do so within reason - no one wants to wait a year for a test to finish.
	return b.Build(
		memstreamer.New(),
		memrecordstore.New(),
		memrolescheduler.New(),
		workflow.WithTimeoutStore(memtimeoutstore.New()),
		workflow.WithClock(clock),
	)
}
