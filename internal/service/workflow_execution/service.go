package workflowexecution

import (
	"context"
	"time"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
)

var _ Service = (*service)(nil)

type Service interface {
	Delete(ctx context.Context, cmd DeleteWorkflowExecutionCommand) error
	SetRunning(ctx context.Context, cmd SetWorkflowExecutionRunningCommand) error
	SetFailed(ctx context.Context, cmd SetWorkflowExecutionFailedCommand) error
	SetCompleted(ctx context.Context, cmd SetWorkflowExecutionCompletedCommand) error

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

func (s service) Delete(ctx context.Context, cmd DeleteWorkflowExecutionCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	return s.workflowExecutionRepo.Delete(ctx, cmd.ID)
}

func (s service) SetRunning(ctx context.Context, cmd SetWorkflowExecutionRunningCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	return s.workflowExecutionRepo.Update(
		ctx,
		cmd.ID,
		func(wfe *model.WorkflowExecution) error {
			wfe.Status = model.WorkflowExecutionStatusRunning
			now := time.Now()
			wfe.StartedAt = &now

			return nil
		},
	)
}

func (s service) SetCompleted(ctx context.Context, cmd SetWorkflowExecutionCompletedCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	return s.workflowExecutionRepo.Update(
		ctx,
		cmd.ID,
		func(wfe *model.WorkflowExecution) error {
			wfe.Status = model.WorkflowExecutionStatusCompleted
			now := time.Now()
			wfe.CompletedAt = &now

			return nil
		},
	)
}

func (s service) SetFailed(ctx context.Context, cmd SetWorkflowExecutionFailedCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	return s.workflowExecutionRepo.Update(
		ctx,
		cmd.ID,
		func(wfe *model.WorkflowExecution) error {
			wfe.Status = model.WorkflowExecutionStatusFailed
			now := time.Now()
			wfe.CompletedAt = &now

			return nil
		},
	)
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
