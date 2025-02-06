package workflow

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/model/workflow/edge"
	"github.com/tuanvumaihuynh/roboflow/internal/model/workflow/node"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerror"
)

type ViewPort struct {
	X    float32 `json:"x" validate:"required"`
	Y    float32 `json:"y" validate:"required"`
	Zoom float32 `json:"zoom" validate:"required"`
}

type Data struct {
	Nodes    []node.Node `json:"nodes" validate:"dive"`
	Edges    []edge.Edge `json:"edges" validate:"dive"`
	Position []float32   `json:"position" validate:"dive"`
	ViewPort ViewPort    `json:"view_port" validate:"required"`
	Zoom     float32     `json:"zoom" validate:"required"`
}

// Validate validates the workflow data.
// It checks if the workflow has at least one node, exactly one trigger node
// and all nodes are connected using DFS.
func (d Data) Validate() error {
	if len(d.Nodes) == 0 {
		return xerror.ValidationFailed(nil, "Workflow must have at least one node")
	}

	// Find trigger node
	var hasTrigger bool
	nodeMap := make(map[string]bool)
	for _, n := range d.Nodes {
		if n.Type == node.TypeTrigger {
			if hasTrigger {
				return xerror.ValidationFailed(nil, "Workflow cannot have multiple trigger nodes")
			}
			hasTrigger = true
		}
		nodeMap[n.ID] = false
	}

	if !hasTrigger {
		return xerror.ValidationFailed(nil, "Workflow must have exactly one trigger node")
	}

	// Validate edges
	for _, edge := range d.Edges {
		if _, exists := nodeMap[edge.Source]; !exists {
			return xerror.ValidationFailed(nil, fmt.Sprintf("Edge source node %s does not exist", edge.Source))
		}
		if _, exists := nodeMap[edge.Target]; !exists {
			return xerror.ValidationFailed(nil, fmt.Sprintf("Edge target node %s does not exist", edge.Target))
		}
	}

	// Check if all nodes are connected using DFS
	visited := make(map[string]bool)
	// dfs is a helper function to check if all nodes are connected using DFS
	var dfs func(nodeID string)
	dfs = func(nodeID string) {
		visited[nodeID] = true
		for _, edge := range d.Edges {
			if edge.Source == nodeID && !visited[edge.Target] {
				dfs(edge.Target)
			}
		}
	}

	// Start DFS from trigger node
	for _, n := range d.Nodes {
		if n.Type == node.TypeTrigger {
			dfs(n.ID)
			break
		}
	}

	// Check if all nodes are visited
	for _, n := range d.Nodes {
		if !visited[n.ID] {
			return xerror.ValidationFailed(nil, fmt.Sprintf("Node %s is not connected to the workflow", n.ID))
		}
	}

	return nil
}

type Workflow struct {
	ID          string
	Name        string
	Description *string
	IsDraft     bool
	IsValid     bool
	Data        Data
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewWorkflow(name string, description *string, isValid bool, data Data) Workflow {
	now := time.Now()
	return Workflow{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		Data:        data,
		IsDraft:     true,
		IsValid:     isValid,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
