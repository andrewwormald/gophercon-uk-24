package main

import (
	"example"

	"k8s.io/utils/clock"

	"github.com/luno/workflow"
)

func main() {
	w := example.NewGopherWorkflow(0, clock.RealClock{})
	err := workflow.MermaidDiagram(w, "./diagram.md", workflow.LeftToRightDirection)
	if err != nil {
		panic(err)
	}
}
