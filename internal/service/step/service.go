package step

import (
	"context"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
)

var _ Service = (*service)(nil)

type Service interface {
	GetByID(ctx context.Context, q GetStepByIDQuery) (model.Step, error)
	List(ctx context.Context, q ListStepQuery) ([]model.Step, error)
}

type service struct {
	stepRepo model.StepRepository
}

func NewService(stepRepo model.StepRepository) *service {
	return &service{
		stepRepo: stepRepo,
	}
}

func (s *service) GetByID(ctx context.Context, q GetStepByIDQuery) (model.Step, error) {
	if err := q.Validate(); err != nil {
		return model.Step{}, err
	}

	return s.stepRepo.Get(ctx, q.ID)
}

func (s *service) List(ctx context.Context, q ListStepQuery) ([]model.Step, error) {
	if err := q.Validate(); err != nil {
		return nil, err
	}

	return s.stepRepo.List(ctx, q.WorkflowExecutionID, q.Sorts)
}
