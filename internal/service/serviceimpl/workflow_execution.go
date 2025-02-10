package serviceimpl

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	stepexecution "github.com/tuanvumaihuynh/roboflow/internal/model/step_execution"
	"github.com/tuanvumaihuynh/roboflow/internal/model/workflow/node"
	workflowexecution "github.com/tuanvumaihuynh/roboflow/internal/model/workflow_execution"
	"github.com/tuanvumaihuynh/roboflow/internal/repository"
	"github.com/tuanvumaihuynh/roboflow/internal/service"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/ptr"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
)

var _ service.WorkflowExecutionService = (*workflowExecutionService)(nil)

type workflowExecutionService struct {
	workflowExecutionRepo repository.WorkflowExecutionRepository
	stepExecutionRepo     repository.StepExecutionRepository
	sqlDBProvider         sqldb.Provider
	validator             validator.Validator
}

func newWorkflowExecutionService(
	workflowExecutionRepo repository.WorkflowExecutionRepository,
	stepExecutionRepo repository.StepExecutionRepository,
	sqlDBProvider sqldb.Provider,
	validator validator.Validator,
) *workflowExecutionService {
	return &workflowExecutionService{
		workflowExecutionRepo: workflowExecutionRepo,
		stepExecutionRepo:     stepExecutionRepo,
		sqlDBProvider:         sqlDBProvider,
		validator:             validator,
	}
}

func (s workflowExecutionService) GetWorkflowExecution(ctx context.Context, params service.GetWorkflowExecutionParams) (workflowexecution.WorkflowExecution, error) {
	if err := s.validator.Validate(params); err != nil {
		return workflowexecution.WorkflowExecution{}, fmt.Errorf("validate params: %w", err)
	}

	we, err := s.workflowExecutionRepo.GetWorkflowExecution(ctx, s.sqlDBProvider.DB(), params.ID)
	if err != nil {
		return workflowexecution.WorkflowExecution{}, fmt.Errorf("repo get workflow execution: %w", err)
	}

	return we, nil
}

func (s workflowExecutionService) ListWorkflowExecutionsByWorkflowID(ctx context.Context, params service.ListWorkflowExecutionsByWorkflowIDParams) (paging.List[workflowexecution.WorkflowExecution], error) {
	if err := s.validator.Validate(params); err != nil {
		return paging.List[workflowexecution.WorkflowExecution]{}, fmt.Errorf("validate params: %w", err)
	}

	we, err := s.workflowExecutionRepo.ListWorkflowExecutionsByWorkflowID(ctx, s.sqlDBProvider.DB(), params.PagingParams, params.Sorts, params.WorkflowID)
	if err != nil {
		return paging.List[workflowexecution.WorkflowExecution]{}, fmt.Errorf("repo list workflow executions by workflow id: %w", err)
	}

	return we, nil
}

func (s workflowExecutionService) ProcessRunWorkflowExecution(ctx context.Context, params service.ProcessRunWorkflowExecutionParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	// Get all steps for this workflow execution
	steps, err := s.stepExecutionRepo.ListStepsByWorkflowExecutionID(
		ctx,
		s.sqlDBProvider.DB(),
		params.WorkflowExecutionID,
	)
	if err != nil {
		return fmt.Errorf("repo list steps by workflow execution id: %w", err)
	}

	// Update workflow execution status to running
	wfe, err := s.workflowExecutionRepo.UpdateWorkflowExecution(
		ctx,
		s.sqlDBProvider.DB(),
		repository.UpdateWorkflowExecutionParams{
			ID:           params.WorkflowExecutionID,
			Status:       workflowexecution.StatusRunning,
			SetStatus:    true,
			StartedAt:    ptr.New(time.Now()),
			SetStartedAt: true,
		},
	)
	if err != nil {
		return fmt.Errorf("repo update workflow execution status: %w", err)
	}

	// Build execution graph
	graph := stepexecution.BuildExecutionGraph(wfe.Data.Edges, steps)

	// Execute workflow
	if err := s.executeWorkflow(ctx, graph); err != nil {
		// Update workflow execution status to failed
		_, err := s.workflowExecutionRepo.UpdateWorkflowExecution(
			ctx,
			s.sqlDBProvider.DB(),
			repository.UpdateWorkflowExecutionParams{
				ID:             params.WorkflowExecutionID,
				Status:         workflowexecution.StatusFailed,
				SetStatus:      true,
				CompletedAt:    ptr.New(time.Now()),
				SetCompletedAt: true,
			},
		)
		if err != nil {
			return fmt.Errorf("repo update workflow execution status: %w", err)
		}
		return fmt.Errorf("execute workflow: %w", err)
	}

	// Update workflow execution status to completed
	_, err = s.workflowExecutionRepo.UpdateWorkflowExecution(
		ctx,
		s.sqlDBProvider.DB(),
		repository.UpdateWorkflowExecutionParams{
			ID:             params.WorkflowExecutionID,
			Status:         workflowexecution.StatusCompleted,
			SetStatus:      true,
			CompletedAt:    ptr.New(time.Now()),
			SetCompletedAt: true,
		},
	)
	if err != nil {
		return fmt.Errorf("repo update workflow execution status: %w", err)
	}

	return nil
}

func (s workflowExecutionService) executeWorkflow(ctx context.Context, graph stepexecution.ExecutionGraph) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(graph))

	// Find trigger node and start execution
	for _, n := range graph {
		if n.Step.Node.Type == node.TypeTrigger {
			wg.Add(1)
			go s.executeNode(ctx, n, &wg, errChan)
			break
		}
	}

	// Wait for all nodes to complete
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collect errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s workflowExecutionService) executeNode(ctx context.Context, n *stepexecution.ExecutionNode, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	// Check if node has already been executed
	n.IsExecutedMu.Lock()
	if n.IsExecuted {
		n.IsExecutedMu.Unlock()
		return
	}
	n.IsExecuted = true
	n.IsExecutedMu.Unlock()

	// Update step status to running
	_, err := s.stepExecutionRepo.UpdateStepExecution(
		ctx,
		s.sqlDBProvider.DB(),
		repository.UpdateStepExecutionParams{
			ID:           n.Step.ID,
			Status:       stepexecution.StatusRunning,
			SetStatus:    true,
			StartedAt:    ptr.New(time.Now()),
			SetStartedAt: true,
		},
	)
	if err != nil {
		errChan <- fmt.Errorf("update step status: %w", err)
	}

	// Execute node logic
	outputs, err := s.executeNodeLogic(ctx, n)
	if err != nil {
		// Update step status to failed
		_, updateErr := s.stepExecutionRepo.UpdateStepExecution(
			ctx,
			s.sqlDBProvider.DB(),
			repository.UpdateStepExecutionParams{
				ID:             n.Step.ID,
				Status:         stepexecution.StatusFailed,
				SetStatus:      true,
				Error:          ptr.New(err.Error()),
				SetError:       true,
				CompletedAt:    ptr.New(time.Now()),
				SetCompletedAt: true,
			},
		)
		if updateErr != nil {
			errChan <- fmt.Errorf("update step status: %w", updateErr)
		}
		errChan <- fmt.Errorf("execute node logic: %w", err)
		return
	}

	// Update step with outputs and completed status
	_, err = s.stepExecutionRepo.UpdateStepExecution(
		ctx,
		s.sqlDBProvider.DB(),
		repository.UpdateStepExecutionParams{
			ID:             n.Step.ID,
			Status:         stepexecution.StatusCompleted,
			SetStatus:      true,
			Outputs:        outputs,
			SetOutputs:     true,
			CompletedAt:    ptr.New(time.Now()),
			SetCompletedAt: true,
		},
	)
	if err != nil {
		errChan <- fmt.Errorf("update step outputs: %w", err)
		return
	}

	// Execute children
	for _, child := range n.Children {
		wg.Add(1)
		go s.executeNode(ctx, child, wg, errChan)
	}
}

// Execute node logic based on node type and return outputs.
func (s workflowExecutionService) executeNodeLogic(_ context.Context, n *stepexecution.ExecutionNode) (map[string]any, error) {
	switch n.Step.Node.Type {
	case node.TypeTrigger:
		return n.Step.Inputs, nil
	// Add other node type handlers here
	default:
		return nil, fmt.Errorf("unsupported node type: %s", n.Step.Node.Type)
	}
}
