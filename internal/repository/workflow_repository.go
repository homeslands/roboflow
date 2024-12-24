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

var _ model.WorkflowRepository = (*workflowRepository)(nil)

type workflowRepository struct {
	store *db.Store
}

func NewWorkflowRepository(store *db.Store) *workflowRepository {
	return &workflowRepository{store: store}
}

func (r workflowRepository) Get(ctx context.Context, id uuid.UUID) (model.Workflow, error) {
	row, err := r.store.GetWorkflow(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Workflow{}, xerrors.ThrowNotFound(err, "workflow not found")
		}
		return model.Workflow{}, err
	}

	return rowWorkflowToModel(*row)
}
func (r workflowRepository) List(ctx context.Context, p paging.Params, sorts []xsort.Sort) (*paging.List[model.Workflow], error) {
	params := db.ListWorkflowsParams{
		Limit:  p.Limit(),
		Offset: p.Offset(),
		Sorts:  sorts,
	}

	row, err := r.store.ListWorkflows(ctx, params)
	if err != nil {
		return nil, err
	}

	var items []model.Workflow
	for _, item := range row.Items {
		m, err := rowWorkflowToModel(item)
		if err != nil {
			return nil, err
		}
		items = append(items, m)
	}

	return paging.NewList(items, row.TotalItem), nil
}
func (r workflowRepository) Create(ctx context.Context, workflow model.Workflow) error {
	def, err := json.Marshal(workflow.Definition)
	if err != nil {
		return xerrors.ThrowInternal(err, "failed to marshal workflow definition")
	}

	err = r.store.CreateWorkflow(ctx, db.CreateWorkflowParams{
		ID:          workflow.ID,
		Name:        workflow.Name,
		Description: workflow.Description,
		Definition:  def,
		CreatedAt:   pgtype.Timestamptz{Time: workflow.CreatedAt, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: workflow.UpdatedAt, Valid: true},
	})

	return err
}
func (r workflowRepository) Update(ctx context.Context, workflow model.Workflow) (model.Workflow, error) {
	def, err := json.Marshal(workflow.Definition)
	if err != nil {
		return model.Workflow{}, xerrors.ThrowInternal(err, "failed to marshal workflow definition")
	}

	wf, err := r.store.UpdateWorkflow(ctx, db.UpdateWorkflowParams{
		ID:          workflow.ID,
		Name:        workflow.Name,
		Description: workflow.Description,
		Definition:  def,
		UpdatedAt:   pgtype.Timestamptz{Time: workflow.UpdatedAt, Valid: true},
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Workflow{}, xerrors.ThrowNotFound(err, "workflow not found")
	}

	return rowWorkflowToModel(*wf)
}
func (r workflowRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteWorkflow(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return xerrors.ThrowNotFound(err, "workflow not found")
	}

	return err
}

func rowWorkflowToModel(row db.Workflow) (model.Workflow, error) {
	m := model.Workflow{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
	if len(row.Definition) > 0 {
		var def model.WorkflowDefinition
		if err := json.Unmarshal(row.Definition, &def); err != nil {
			return model.Workflow{}, xerrors.ThrowInternal(err, "failed to unmarshal workflow definition")
		}
		m.Definition = &def
	}

	return m, nil
}
