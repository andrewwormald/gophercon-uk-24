package example_test

import (
	"context"
	"testing"
	"time"

	"github.com/luno/workflow"
	clock_testing "k8s.io/utils/clock/testing"

	"example"
)

func TestNewGopherWorkflow(t *testing.T) {
	now := time.Now()
	clock := clock_testing.NewFakeClock(now)
	age := 12
	w := example.NewGopherWorkflow(age, clock)

	ctx := context.Background()
	w.Run(ctx)
	defer w.Stop()

	// foreignID can be anything as long as it references something else in the system and is
	// not complete silliness.
	foreignID := "1"
	runID, err := w.Trigger(ctx, foreignID, example.StatusStarted)
	if err != nil {
		panic(err)
	}

	// Wait for the timeout to be inserted before updating the clock. Alternatively a record
	// would be stored in the TimeoutStore with the updated time and still wait until the
	// timeout expires which may never happen if using a fake clock.
	workflow.AwaitTimeoutInsert(t, w, foreignID, runID, example.StatusSentToSchool)

	clock.SetTime(now.AddDate(100, 0, 0))

	gr := example.GraduationResponse{
		Graduated: true,
	}

	// Calling TriggerCallbackOn allows for simulating a callback as soon as the workflow gets to
	// the provided status. Outside of testing, you would need to use the workflow to trigger a
	// callback which looks like this: w.Callback(ctx, foreignID, example.StatusFinishedSchool, {{io.Reader}})
	workflow.TriggerCallbackOn(t, w, foreignID, runID, example.StatusFinishedSchool, gr)

	// Require waits for the provided status to be reached and requires that the end result matches the
	// provided result.
	workflow.Require(t, w, foreignID, example.StatusSentToWork, example.Gopher{
		Name:   "Budgie",
		Colour: "Blue",
		Age:    12,
		School: "Graduated school and ready to find lost goroutines",
		Work:   "Sent to start finding lost goroutines",
	})
}
