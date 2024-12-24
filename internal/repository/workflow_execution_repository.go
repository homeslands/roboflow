package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/tuanvumaihuynh/roboflow/db"
	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

var _ model.WorkflowExecutionRepository = (*workflowExecutionRepository)(nil)

type workflowExecutionRepository struct {
	store *db.Store
}

func NewWorkflowExecutionRepository(store *db.Store) *workflowExecutionRepository {
	return &workflowExecutionRepository{store: store}
}

func (r workflowExecutionRepository) Get(ctx context.Context, id uuid.UUID) (model.WorkflowExecution, error) {
	row, err := r.store.GetWorkflowExecution(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.WorkflowExecution{}, xerrors.ThrowNotFound(err, "workflow execution not found")
		}
		return model.WorkflowExecution{}, err
	}

	return rowWorkflowExecutionToModel(*row)
}

func (r workflowExecutionRepository) GetStatus(ctx context.Context, id uuid.UUID) (model.WorkflowExecutionStatus, error) {
	row, err := r.store.GetWorkflowExecutionStatus(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", xerrors.ThrowNotFound(err, "workflow execution not found")
		}
		return "", err
	}

	return model.WorkflowExecutionStatus(row), nil
}

func (r workflowExecutionRepository) List(ctx context.Context, workflowID uuid.UUID, p paging.Params, sorts []xsort.Sort) (*paging.List[model.WorkflowExecution], error) {
	// TODO: Implement this method
	return nil, nil
}
func (r workflowExecutionRepository) Create(ctx context.Context, workflowExecution model.WorkflowExecution) error {
	def, err := json.Marshal(workflowExecution.Definition)
	if err != nil {
		return xerrors.ThrowInternal(err, "failed to marshal workflow definition")
	}

	createWfeParams := db.CreateWorkflowExecutionParams{
		ID:         workflowExecution.ID,
		WorkflowID: workflowExecution.WorkflowID,
		Status:     string(workflowExecution.Status),
		Env:        workflowExecution.Env,
		Definition: def,
		CreatedAt:  pgtype.Timestamptz{Time: workflowExecution.CreatedAt, Valid: true},
	}

	createStepsParams := make([]db.BulkInsertStepsParams, 0, len(*workflowExecution.Steps))
	for _, step := range *workflowExecution.Steps {
		node, err := json.Marshal(step.Node)
		if err != nil {
			return xerrors.ThrowInternal(err, "failed to marshal step node")
		}

		param := db.BulkInsertStepsParams{
			ID:                  step.ID,
			WorkflowExecutionID: step.WorkflowExecutionID,
			Env:                 step.Env,
			Node:                node,
			Status:              string(step.Status),
		}
		createStepsParams = append(createStepsParams, param)
	}

	err = r.store.WithTx(ctx, func(s db.Store) error {
		err := s.CreateWorkflowExecution(ctx, createWfeParams)
		if err != nil {
			return err
		}

		_, err = s.BulkInsertSteps(ctx, createStepsParams)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return xerrors.ThrowInternal(err, "failed to create workflow execution")
	}

	return nil
}
func (r workflowExecutionRepository) Update(ctx context.Context, workflowExecution model.WorkflowExecution) (model.WorkflowExecution, error) {
	params := db.UpdateWorkflowExecutionParams{
		ID:     workflowExecution.ID,
		Status: string(workflowExecution.Status),
	}
	if workflowExecution.StartedAt != nil {
		params.StartedAt = pgtype.Timestamptz{Time: *workflowExecution.StartedAt, Valid: true}
	}
	if workflowExecution.CompletedAt != nil {
		params.CompletedAt = pgtype.Timestamptz{Time: *workflowExecution.CompletedAt, Valid: true}
	}

	row, err := r.store.UpdateWorkflowExecution(ctx, params)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.WorkflowExecution{}, xerrors.ThrowNotFound(err, "workflow execution not found")
	}

	return rowWorkflowExecutionToModel(*row)
}
func (r workflowExecutionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteWorkflowExecution(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return xerrors.ThrowNotFound(err, "workflow execution not found")
	}

	return err
}

func rowWorkflowExecutionToModel(row db.WorkflowExecution) (model.WorkflowExecution, error) {
	m := model.WorkflowExecution{
		ID:         row.ID,
		WorkflowID: row.WorkflowID,
		Status:     model.WorkflowExecutionStatus(row.Status),
		Env:        row.Env,
		CreatedAt:  row.CreatedAt.Time,
	}
	if len(row.Definition) > 0 {
		var def model.WorkflowDefinition
		if err := json.Unmarshal(row.Definition, &def); err != nil {
			return model.WorkflowExecution{}, xerrors.ThrowInternal(err, "failed to unmarshal workflow definition")
		}
		m.Definition = def
	}

	if row.StartedAt.Valid {
		m.StartedAt = &row.StartedAt.Time
	}
	if row.CompletedAt.Valid {
		m.CompletedAt = &row.CompletedAt.Time
	}

	return m, nil
}
