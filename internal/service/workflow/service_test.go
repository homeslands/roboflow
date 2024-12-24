package workflow_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/internal/model/mocks"
	"github.com/tuanvumaihuynh/roboflow/internal/service/workflow"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	pubsubmocks "github.com/tuanvumaihuynh/roboflow/pkg/pubsub/mocks"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

func TestWorkflowService(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	publisher := pubsubmocks.NewFakePublisher(t)
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		testCases := []struct {
			name         string
			cmd          workflow.CreateWorkflowCommand
			mockBehavior func(*mocks.FakeWorkflowRepository)
			shouldErr    bool
		}{
			{
				name: "Should create successfully",
				cmd: workflow.CreateWorkflowCommand{
					Name:        "Test Workflow",
					Description: nil,
					Definition:  validWorkflowDef,
				},
				mockBehavior: func(r *mocks.FakeWorkflowRepository) {
					r.On("Create", ctx, mock.Anything).Return(nil)
				},
				shouldErr: false,
			},
			{
				name: "Should validate command before do anything",
				cmd: workflow.CreateWorkflowCommand{
					Name:        "",
					Description: nil,
					Definition:  validWorkflowDef,
				},
				mockBehavior: func(r *mocks.FakeWorkflowRepository) {
					r.AssertNotCalled(t, "Create", ctx, mock.Anything)
				},
				shouldErr: true,
			},
			{
				name: "Should return error if repository return error",
				cmd: workflow.CreateWorkflowCommand{
					Name:        "Test Workflow",
					Description: nil,
					Definition:  validWorkflowDef,
				},
				mockBehavior: func(r *mocks.FakeWorkflowRepository) {
					r.On("Create", ctx, mock.Anything).Return(assert.AnError)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				wfRepo := mocks.NewFakeWorkflowRepository(t)
				tc.mockBehavior(wfRepo)
				wfeRepo := mocks.NewFakeWorkflowExecutionRepository(t)

				s := workflow.NewService(wfRepo, wfeRepo, publisher, log)
				result, err := s.Create(ctx, tc.cmd)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.NotEmpty(t, result.ID)
					assert.Equal(t, tc.cmd.Name, result.Name)
					assert.Equal(t, tc.cmd.Description, result.Description)
					assert.Equal(t, tc.cmd.Definition, *result.Definition)
					assert.NotEmpty(t, result.CreatedAt)
					assert.NotEmpty(t, result.UpdatedAt)
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("Update", func(t *testing.T) {
		testCases := []struct {
			name         string
			cmd          workflow.UpdateWorkflowCommand
			mockBehavior func(*mocks.FakeWorkflowRepository)
			shouldErr    bool
		}{
			{
				name: "Should update successfully",
				cmd: workflow.UpdateWorkflowCommand{
					ID:          validID,
					Name:        "Test Workflow",
					Description: nil,
					Definition:  validWorkflowDef,
				},
				mockBehavior: func(r *mocks.FakeWorkflowRepository) {
					r.On("Update", ctx, mock.Anything).Return(validWorkflow, nil)
				},
				shouldErr: false,
			},
			{
				name: "Should validate command before do anything",
				cmd: workflow.UpdateWorkflowCommand{
					ID:          validID,
					Name:        "",
					Description: nil,
					Definition:  validWorkflowDef,
				},
				mockBehavior: func(r *mocks.FakeWorkflowRepository) {
					r.AssertNotCalled(t, "Update", ctx, mock.Anything)
				},
				shouldErr: true,
			},
			{
				name: "Should return error if repository return error",
				cmd: workflow.UpdateWorkflowCommand{
					ID:          validID,
					Name:        "Test Workflow",
					Description: nil,
					Definition:  validWorkflowDef,
				},
				mockBehavior: func(r *mocks.FakeWorkflowRepository) {
					r.On("Update", ctx, mock.Anything).Return(validWorkflow, assert.AnError)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				wfRepo := mocks.NewFakeWorkflowRepository(t)
				tc.mockBehavior(wfRepo)
				wfeRepo := mocks.NewFakeWorkflowExecutionRepository(t)

				s := workflow.NewService(wfRepo, wfeRepo, publisher, log)
				result, err := s.Update(ctx, tc.cmd)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.NotEmpty(t, result.ID)
					assert.Equal(t, tc.cmd.Name, result.Name)
					assert.Equal(t, tc.cmd.Description, result.Description)
					assert.Equal(t, tc.cmd.Definition, *result.Definition)
					assert.NotEmpty(t, result.CreatedAt)
					assert.NotEmpty(t, result.UpdatedAt)
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("Delete", func(t *testing.T) {
		testCases := []struct {
			name         string
			cmd          workflow.DeleteWorkflowCommand
			mockBehavior func(*mocks.FakeWorkflowRepository)
			shouldErr    bool
		}{
			{
				name: "Should delete successfully",
				cmd: workflow.DeleteWorkflowCommand{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeWorkflowRepository) {
					r.On("Delete", ctx, validID).Return(nil)
				},
				shouldErr: false,
			},
			{
				name: "Should validate command before do anything",
				cmd: workflow.DeleteWorkflowCommand{
					ID: uuid.Nil,
				},
				mockBehavior: func(r *mocks.FakeWorkflowRepository) {
					r.AssertNotCalled(t, "Delete", ctx, validID)
				},
				shouldErr: true,
			},
			{
				name: "Should return error if repository return error",
				cmd: workflow.DeleteWorkflowCommand{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeWorkflowRepository) {
					r.On("Delete", ctx, validID).Return(assert.AnError)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				wfRepo := mocks.NewFakeWorkflowRepository(t)
				tc.mockBehavior(wfRepo)
				wfeRepo := mocks.NewFakeWorkflowExecutionRepository(t)

				s := workflow.NewService(wfRepo, wfeRepo, publisher, log)
				err := s.Delete(ctx, tc.cmd)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("Run", func(t *testing.T) {})
	t.Run("GetByID", func(t *testing.T) {
		testCases := []struct {
			name         string
			query        workflow.GetWorkflowByIDQuery
			mockBehavior func(*mocks.FakeWorkflowRepository)
			shouldErr    bool
		}{
			{
				name: "Should get workflow successfully",
				query: workflow.GetWorkflowByIDQuery{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeWorkflowRepository) {
					r.On("Get", ctx, validID).Return(validWorkflow, nil)
				},
			},
			{
				name: "Should validate query before do anything",
				query: workflow.GetWorkflowByIDQuery{
					ID: uuid.Nil,
				},
				mockBehavior: func(r *mocks.FakeWorkflowRepository) {
					r.AssertNotCalled(t, "Get", ctx, validID)
				},
				shouldErr: true,
			},
			{
				name: "Should return error if repository return error",
				query: workflow.GetWorkflowByIDQuery{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeWorkflowRepository) {
					r.On("Get", ctx, validID).Return(model.Workflow{}, assert.AnError)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				wfRepo := mocks.NewFakeWorkflowRepository(t)
				tc.mockBehavior(wfRepo)
				wfeRepo := mocks.NewFakeWorkflowExecutionRepository(t)

				s := workflow.NewService(wfRepo, wfeRepo, publisher, log)
				result, err := s.GetByID(ctx, tc.query)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.NotEmpty(t, result.ID)
					assert.Equal(t, validWorkflow.Name, result.Name)
					assert.Equal(t, validWorkflow.Description, result.Description)
					assert.Equal(t, validWorkflow.Definition, result.Definition)
					assert.NotEmpty(t, result.CreatedAt)
					assert.NotEmpty(t, result.UpdatedAt)
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("List", func(t *testing.T) {
		validPagingWorkflowSteps := paging.List[model.Workflow]{
			Items:     []model.Workflow{validWorkflow},
			TotalItem: 1,
		}

		testCases := []struct {
			name         string
			query        workflow.ListWorkflowQuery
			mockBehavior func(*mocks.FakeWorkflowRepository)
			shouldErr    bool
		}{
			{
				name: "Should list workflow successfully",
				query: workflow.ListWorkflowQuery{
					PagingParams: paging.NewParams(nil, nil),
					Sorts:        []xsort.Sort{},
				},
				mockBehavior: func(r *mocks.FakeWorkflowRepository) {
					r.On("List", ctx, mock.Anything, mock.Anything).Return(&validPagingWorkflowSteps, nil)
				},
				shouldErr: false,
			},
			{
				name: "Should validate query before do anything",
				query: workflow.ListWorkflowQuery{
					PagingParams: paging.NewParams(nil, nil),
					Sorts: []xsort.Sort{
						{
							Col:   "invalid",
							Order: xsort.OrderASC,
						},
					},
				},
				mockBehavior: func(r *mocks.FakeWorkflowRepository) {
					r.AssertNotCalled(t, "List", ctx, mock.Anything, mock.Anything)
				},
				shouldErr: true,
			},
			{
				name: "Should return error if repository return error",
				query: workflow.ListWorkflowQuery{
					PagingParams: paging.NewParams(nil, nil),
					Sorts:        []xsort.Sort{},
				},
				mockBehavior: func(r *mocks.FakeWorkflowRepository) {
					r.On("List", ctx, mock.Anything, mock.Anything).Return(nil, assert.AnError)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				wfRepo := mocks.NewFakeWorkflowRepository(t)
				tc.mockBehavior(wfRepo)
				wfeRepo := mocks.NewFakeWorkflowExecutionRepository(t)

				s := workflow.NewService(wfRepo, wfeRepo, publisher, log)
				result, err := s.List(ctx, tc.query)

				if tc.shouldErr {
					assert.Nil(t, result)
					assert.Error(t, err)
				} else {
					assert.NotEmpty(t, result.Items)
					assert.NoError(t, err)
				}
			})
		}
	})
}

var (
	validID          = uuid.New()
	validWorkflowDef = model.WorkflowDefinition{
		Nodes: []model.WorkflowNode{
			{
				ID:          "1",
				Type:        model.NodeTypeRaybot,
				Initialized: true,
				Position: struct {
					X float32 `json:"x" validate:"required"`
					Y float32 `json:"y" validate:"required"`
				}{
					X: 1,
					Y: 1,
				},
				Definition: model.NodeDefinition{
					Type:   model.TaskRaybotCheckQrCode,
					Fields: map[string]model.NodeField{},
				},
			},
		},
		Edges: []model.WorkflowEdge{
			{
				ID:           "1",
				Type:         "type",
				Source:       "1",
				Target:       "2",
				SourceHandle: "source",
				TargetHandle: "target",
				Label:        "label",
				Animated:     true,
				SourceX:      1.0,
				SourceY:      1.0,
				TargetX:      1.0,
				TargetY:      1.0,
			},
		},
		Position: []float32{1.0},
		ViewPort: model.ViewPort{
			X:    float32(1.0),
			Y:    float32(1.0),
			Zoom: float32(1.0),
		},
		Zoom: float32(1.0),
	}
	validWorkflow = model.Workflow{
		ID:          validID,
		Name:        "Test Workflow",
		Description: nil,
		Definition:  &validWorkflowDef,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
)
