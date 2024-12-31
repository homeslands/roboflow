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
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

var _ model.StepRepository = (*StepRepository)(nil)

type StepRepository struct {
	store *db.Store
}

func NewStepRepository(store *db.Store) *StepRepository {
	return &StepRepository{store: store}
}

func (r StepRepository) Get(ctx context.Context, id uuid.UUID) (model.Step, error) {
	row, err := r.store.GetStepByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Step{}, xerrors.ThrowNotFound(nil, "step not found")
		}
	}

	return rowStepToModel(*row)
}
func (r StepRepository) List(ctx context.Context, workflowExecutionID uuid.UUID, sorts []xsort.Sort) ([]model.Step, error) {
	rows, err := r.store.ListSteps(ctx, db.ListStepsParams{
		WorkflowExecutionID: workflowExecutionID,
		Sorts:               sorts,
	})
	if err != nil {
		return nil, xerrors.ThrowInternal(err, "failed to list steps")
	}

	var items []model.Step
	for _, row := range rows {
		m, err := rowStepToModel(row)
		if err != nil {
			return nil, xerrors.ThrowInternal(err, "failed to convert row to model")
		}
		items = append(items, m)
	}

	return items, nil
}

func (r StepRepository) Update(ctx context.Context, step model.Step) error {
	param := db.UpdateStepParams{
		ID:                  step.ID,
		WorkflowExecutionID: step.WorkflowExecutionID,
		Env:                 step.Env,
		Status:              string(step.Status),
	}
	nodeBytes, err := json.Marshal(step.Node)
	if err != nil {
		return xerrors.ThrowInternal(err, "failed to marshal node")
	}
	param.Node = nodeBytes
	if step.StartedAt != nil {
		param.StartedAt = pgtype.Timestamptz{Time: *step.StartedAt, Valid: true}
	}
	if step.CompletedAt != nil {
		param.CompletedAt = pgtype.Timestamptz{Time: *step.CompletedAt, Valid: true}
	}

	if err := r.store.UpdateStep(ctx, param); err != nil {
		return xerrors.ThrowInternal(err, "failed to update step")
	}

	return nil
}

func rowStepToModel(row db.Step) (model.Step, error) {
	m := model.Step{
		ID:                  row.ID,
		WorkflowExecutionID: row.WorkflowExecutionID,
		Env:                 row.Env,
		Status:              model.WorkflowExecutionStepStatus(row.Status),
	}

	if len(row.Node) > 0 {
		var node model.WorkflowNode
		if err := json.Unmarshal([]byte(row.Node), &node); err != nil {
			return model.Step{}, xerrors.ThrowInternal(err, "failed to unmarshal node")
		}
		m.Node = node
	}

	if row.StartedAt.Valid {
		m.StartedAt = &row.StartedAt.Time
	}
	if row.CompletedAt.Valid {
		m.CompletedAt = &row.CompletedAt.Time
	}

	return m, nil
}
