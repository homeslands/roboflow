package stepexecution

import (
	"sync"

	"github.com/tuanvumaihuynh/roboflow/internal/model/workflow/edge"
)

// ExecutionNode is a node in the execution graph
type ExecutionNode struct {
	Step     StepExecution
	Parents  []*ExecutionNode
	Children []*ExecutionNode
	Outputs  map[string]any

	IsExecuted   bool
	IsExecutedMu sync.Mutex
}

// ExecutionGraph is a map of node ID to execution node
type ExecutionGraph map[string]*ExecutionNode

// BuildExecutionGraph builds an execution graph from a list of edges and steps
func BuildExecutionGraph(edges []edge.Edge, steps []StepExecution) ExecutionGraph {
	nodes := make(ExecutionGraph)

	// Create nodes
	for _, step := range steps {
		nodes[step.Node.ID] = &ExecutionNode{
			Step:    step,
			Outputs: make(map[string]any),
		}
	}

	// Connect nodes based on edges
	for _, edge := range edges {
		// Get parent and child nodes
		if parent, ok := nodes[edge.Source]; ok {
			if child, ok := nodes[edge.Target]; ok {
				parent.Children = append(parent.Children, child)
				child.Parents = append(child.Parents, parent)
			}
		}
	}

	return nodes
}
