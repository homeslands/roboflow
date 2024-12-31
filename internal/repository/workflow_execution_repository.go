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
func (r workflowExecutionRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	fn func(wfe *model.WorkflowExecution) error,
) error {
	return r.store.WithTx(ctx, func(s db.Store) error {
		wfeRow, err := s.GetWorkflowExecutionForUpdate(ctx, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return xerrors.ThrowNotFound(err, "workflow execution not found")
			}
			return err
		}

		wfe, err := rowWorkflowExecutionToModel(*wfeRow)
		if err != nil {
			return err
		}

		if err := fn(&wfe); err != nil {
			return err
		}

		def, err := json.Marshal(wfe.Definition)
		if err != nil {
			return xerrors.ThrowInternal(err, "failed to marshal workflow definition")
		}

		updateWfeParams := db.UpdateWorkflowExecutionParams{
			ID:         wfe.ID,
			WorkflowID: wfe.WorkflowID,
			Status:     string(wfe.Status),
			Env:        wfe.Env,
			Definition: def,
			CreatedAt:  pgtype.Timestamptz{Time: wfe.CreatedAt, Valid: true},
		}
		if wfe.StartedAt != nil {
			updateWfeParams.StartedAt = pgtype.Timestamptz{Time: *wfe.StartedAt, Valid: true}
		}
		if wfe.CompletedAt != nil {
			updateWfeParams.CompletedAt = pgtype.Timestamptz{Time: *wfe.CompletedAt, Valid: true}
		}

		err = s.UpdateWorkflowExecution(ctx, updateWfeParams)
		if err != nil {
			return err
		}

		return nil
	})
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
