package main

import (
	"example"
	"github.com/luno/workflow"
	"k8s.io/utils/clock"
)

func main() {
	w := example.NewGopherWorkflow(0, clock.RealClock{})
	err := workflow.MermaidDiagram(w, "./diagram.md", workflow.LeftToRightDirection)
	if err != nil {
		panic(err)
	}
}
