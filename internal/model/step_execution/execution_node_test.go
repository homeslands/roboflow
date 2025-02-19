package stepexecution

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tuanvumaihuynh/roboflow/internal/model/workflow/edge"
	"github.com/tuanvumaihuynh/roboflow/internal/model/workflow/node"
)

func TestBuildExecutionGraph(t *testing.T) {
	// Setup test data
	steps := []StepExecution{
		{Node: node.Node{ID: "1"}},
		{Node: node.Node{ID: "2"}},
		{Node: node.Node{ID: "3"}},
	}

	edges := []edge.Edge{
		{Source: "1", Target: "2"},
		{Source: "2", Target: "3"},
	}

	// Execute
	graph := BuildExecutionGraph(edges, steps)

	// Assert
	assert.Len(t, graph, 3)

	// Check node 1
	node1 := graph["1"]
	assert.Equal(t, "1", node1.Step.Node.ID)
	assert.Empty(t, node1.Parents)
	assert.Len(t, node1.Children, 1)
	assert.Equal(t, "2", node1.Children[0].Step.Node.ID)

	// Check node 2
	node2 := graph["2"]
	assert.Equal(t, "2", node2.Step.Node.ID)
	assert.Len(t, node2.Parents, 1)
	assert.Equal(t, "1", node2.Parents[0].Step.Node.ID)
	assert.Len(t, node2.Children, 1)
	assert.Equal(t, "3", node2.Children[0].Step.Node.ID)

	// Check node 3
	node3 := graph["3"]
	assert.Equal(t, "3", node3.Step.Node.ID)
	assert.Len(t, node3.Parents, 1)
	assert.Equal(t, "2", node3.Parents[0].Step.Node.ID)
	assert.Empty(t, node3.Children)
}
