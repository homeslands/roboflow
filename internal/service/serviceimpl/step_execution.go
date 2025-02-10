package serviceimpl

import (
	"context"
	"fmt"

	"github.com/tuanvumaihuynh/roboflow/internal/db/sqldb"
	stepexecution "github.com/tuanvumaihuynh/roboflow/internal/model/step_execution"
	"github.com/tuanvumaihuynh/roboflow/internal/repository"
	"github.com/tuanvumaihuynh/roboflow/internal/service"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
)

var _ service.StepExecutionService = (*stepExecutionService)(nil)

type stepExecutionService struct {
	stepExecutionRepo repository.StepExecutionRepository
	sqlDBProvider     sqldb.Provider
	validator         validator.Validator
}

func newStepExecutionService(
	stepExecutionRepo repository.StepExecutionRepository,
	sqlDBProvider sqldb.Provider,
	validator validator.Validator,
) *stepExecutionService {
	return &stepExecutionService{
		stepExecutionRepo: stepExecutionRepo,
		sqlDBProvider:     sqlDBProvider,
		validator:         validator,
	}
}

func (s stepExecutionService) GetStepExecution(ctx context.Context, params service.GetStepExecutionParams) (stepexecution.StepExecution, error) {
	if err := s.validator.Validate(params); err != nil {
		return stepexecution.StepExecution{}, fmt.Errorf("validate params: %w", err)
	}

	stepExecution, err := s.stepExecutionRepo.GetStepExecution(ctx, s.sqlDBProvider.DB(), params.ID)
	if err != nil {
		return stepexecution.StepExecution{}, fmt.Errorf("repo get step execution: %w", err)
	}

	return stepExecution, nil
}

func (s stepExecutionService) ListStepsByWorkflowExecutionID(ctx context.Context, params service.ListStepsByWorkflowExecutionIDParams) ([]stepexecution.StepExecution, error) {
	if err := s.validator.Validate(params); err != nil {
		return nil, fmt.Errorf("validate params: %w", err)
	}

	steps, err := s.stepExecutionRepo.ListStepsByWorkflowExecutionID(ctx, s.sqlDBProvider.DB(), params.WorkflowExecutionID)
	if err != nil {
		return nil, fmt.Errorf("repo list steps by workflow execution id: %w", err)
	}

	return steps, nil
}

func (s stepExecutionService) UpdateStepExecution(ctx context.Context, params service.UpdateStepExecutionParams) (stepexecution.StepExecution, error) {
	if err := s.validator.Validate(params); err != nil {
		return stepexecution.StepExecution{}, fmt.Errorf("validate params: %w", err)
	}

	ret, err := s.stepExecutionRepo.UpdateStepExecution(
		ctx,
		s.sqlDBProvider.DB(),
		repository.UpdateStepExecutionParams{
			ID:             params.ID,
			Status:         params.Status,
			SetStatus:      params.SetStatus,
			Inputs:         params.Inputs,
			SetInputs:      params.SetInputs,
			Outputs:        params.Outputs,
			SetOutputs:     params.SetOutputs,
			Error:          params.Error,
			SetError:       params.SetError,
			StartedAt:      params.StartedAt,
			SetStartedAt:   params.SetStartedAt,
			CompletedAt:    params.CompletedAt,
			SetCompletedAt: params.SetCompletedAt,
		},
	)
	if err != nil {
		return stepexecution.StepExecution{}, fmt.Errorf("repo update step execution: %w", err)
	}

	return ret, nil
}
