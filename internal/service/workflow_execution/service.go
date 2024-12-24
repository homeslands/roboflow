package workflowexecution

import (
	"context"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
)

var _ Service = (*service)(nil)

type Service interface {
	Update(ctx context.Context, cmd UpdateWorkflowExecutionCommand) (model.WorkflowExecution, error)
	Delete(ctx context.Context, cmd DeleteWorkflowExecutionCommand) error

	GetByID(ctx context.Context, q GetWorkflowExecutionByIDQuery) (model.WorkflowExecution, error)
	GetStatusByID(ctx context.Context, q GetWorkflowExecutionStatusByIDQuery) (model.WorkflowExecutionStatus, error)
	List(ctx context.Context, q ListWorkflowExecutionQuery) (*paging.List[model.WorkflowExecution], error)
}

type service struct {
	workflowExecutionRepo model.WorkflowExecutionRepository
}

func NewService(workflowExecutionRepo model.WorkflowExecutionRepository) *service {
	return &service{
		workflowExecutionRepo: workflowExecutionRepo,
	}
}

func (s service) Update(ctx context.Context, cmd UpdateWorkflowExecutionCommand) (model.WorkflowExecution, error) {
	if err := cmd.Validate(); err != nil {
		return model.WorkflowExecution{}, err
	}

	w := model.WorkflowExecution{
		ID:          cmd.ID,
		Status:      cmd.Status,
		StartedAt:   cmd.StartedAt,
		CompletedAt: cmd.CompletedAt,
	}

	exec, err := s.workflowExecutionRepo.Update(ctx, w)
	if err != nil {
		return model.WorkflowExecution{}, err
	}

	return exec, nil
}

func (s service) Delete(ctx context.Context, cmd DeleteWorkflowExecutionCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	return s.workflowExecutionRepo.Delete(ctx, cmd.ID)
}

func (s service) GetByID(ctx context.Context, q GetWorkflowExecutionByIDQuery) (model.WorkflowExecution, error) {
	if err := q.Validate(); err != nil {
		return model.WorkflowExecution{}, err
	}

	return s.workflowExecutionRepo.Get(ctx, q.ID)
}

func (s service) GetStatusByID(ctx context.Context, q GetWorkflowExecutionStatusByIDQuery) (model.WorkflowExecutionStatus, error) {
	if err := q.Validate(); err != nil {
		return "", err
	}

	return s.workflowExecutionRepo.GetStatus(ctx, q.ID)
}

func (s service) List(ctx context.Context, q ListWorkflowExecutionQuery) (*paging.List[model.WorkflowExecution], error) {
	if err := q.Validate(); err != nil {
		return nil, err
	}

	return s.workflowExecutionRepo.List(ctx, q.WorkflowID, q.PagingParams, q.Sorts)
}
