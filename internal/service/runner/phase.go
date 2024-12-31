package runner

import (
	"container/list"
	"errors"
	"fmt"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
)

type Phase struct {
	Steps []model.Step
}

func generatePhases(wf model.WorkflowExecution) ([]Phase, error) {
	if wf.Steps == nil {
		return nil, errors.New("steps not found")
	}

	nodeMap := make(map[string]*model.WorkflowNode)
	incomingEdges := make(map[string][]string)
	outgoingEdges := make(map[string][]string)

	// Populate node and edge relationship maps
	for _, node := range wf.Definition.Nodes {
		nodeMap[node.ID] = &node
	}
	for _, edge := range wf.Definition.Edges {
		incomingEdges[edge.Target] = append(incomingEdges[edge.Target], edge.Source)
		outgoingEdges[edge.Source] = append(outgoingEdges[edge.Source], edge.Target)
	}

	// Find the initial step
	initStep, err := findInitStep(*wf.Steps)
	if err != nil {
		return nil, fmt.Errorf("failed to find init step: %w", err)
	}

	phases := []Phase{}
	visitedNodes := make(map[string]bool)
	currentPhase := Phase{Steps: []model.Step{*initStep}}
	visitedNodes[initStep.Node.ID] = true

	// Queue to manage nodes to process
	queue := list.New()
	for _, outgoingNode := range outgoingEdges[initStep.Node.ID] {
		queue.PushBack(outgoingNode)
	}

	// Process nodes to create phases
	for queue.Len() > 0 {
		currentNodeID := queue.Remove(queue.Front()).(string)

		// Skip if node already visited
		if visitedNodes[currentNodeID] {
			continue
		}

		// Check if all incoming edges's source nodes are visited
		canProcess := true
		for _, incomingNodeID := range incomingEdges[currentNodeID] {
			if !visitedNodes[incomingNodeID] {
				canProcess = false
				break
			}
		}

		if !canProcess {
			// Not ready to process, put back in queue
			queue.PushBack(currentNodeID)
			continue
		}

		// Find corresponding step for the node
		var currentStep model.Step
		for _, step := range *wf.Steps {
			if step.Node.ID == currentNodeID {
				currentStep = step
				break
			}
		}

		// Add step to current phase or create new phase
		if len(currentPhase.Steps) == 0 {
			currentPhase.Steps = []model.Step{currentStep}
		} else {
			phases = append(phases, currentPhase)
			currentPhase = Phase{Steps: []model.Step{currentStep}}
		}

		visitedNodes[currentNodeID] = true

		// Add next nodes to queue
		for _, nextNodeID := range outgoingEdges[currentNodeID] {
			if !visitedNodes[nextNodeID] {
				queue.PushBack(nextNodeID)
			}
		}
	}

	// Add last phase if not empty
	if len(currentPhase.Steps) > 0 {
		phases = append(phases, currentPhase)
	}

	return phases, nil
}

func findInitStep(steps []model.Step) (*model.Step, error) {
	for _, step := range steps {
		if step.Node.Definition.Type == model.TaskRaybotValidateState {
			return &step, nil
		}
	}

	return nil, errors.New("init step not found")
}
