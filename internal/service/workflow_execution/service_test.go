package workflowexecution_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/internal/model/mocks"
	workflowexecution "github.com/tuanvumaihuynh/roboflow/internal/service/workflow_execution"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

func TestWorkflowExecutionService(t *testing.T) {
	ctx := context.Background()
	t.Run("Update", func(t *testing.T) {
		testCases := []struct {
			name         string
			cmd          workflowexecution.UpdateWorkflowExecutionCommand
			mockBehavior func(*mocks.FakeWorkflowExecutionRepository)
			shouldErr    bool
		}{
			{
				name: "Should update successfully",
				cmd: workflowexecution.UpdateWorkflowExecutionCommand{
					ID:          validID,
					Status:      model.WorkflowExecutionStatusCompleted,
					StartedAt:   nil,
					CompletedAt: nil,
				},
				mockBehavior: func(r *mocks.FakeWorkflowExecutionRepository) {
					r.On("Update", ctx, mock.Anything).Return(validWorkflowExecution, nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when validate command failed",
				cmd: workflowexecution.UpdateWorkflowExecutionCommand{
					ID:          uuid.Nil,
					Status:      model.WorkflowExecutionStatusCompleted,
					StartedAt:   nil,
					CompletedAt: nil,
				},
				mockBehavior: func(r *mocks.FakeWorkflowExecutionRepository) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository fails to update",
				cmd: workflowexecution.UpdateWorkflowExecutionCommand{
					ID:          validID,
					Status:      model.WorkflowExecutionStatusCompleted,
					StartedAt:   nil,
					CompletedAt: nil,
				},
				mockBehavior: func(r *mocks.FakeWorkflowExecutionRepository) {
					r.On("Update", ctx, mock.Anything).Return(model.WorkflowExecution{}, assert.AnError)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				repo := mocks.NewFakeWorkflowExecutionRepository(t)
				tc.mockBehavior(repo)

				s := workflowexecution.NewService(repo)
				result, err := s.Update(ctx, tc.cmd)

				if tc.shouldErr {
					// assert.Nil(t, err)
					assert.Error(t, err)
				} else {
					assert.Equal(t, tc.cmd.ID, result.ID)
					assert.Equal(t, tc.cmd.Status, result.Status)
					assert.Equal(t, tc.cmd.StartedAt, result.StartedAt)
					assert.Equal(t, tc.cmd.CompletedAt, result.CompletedAt)
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("Delete", func(t *testing.T) {
		testCases := []struct {
			name         string
			cmd          workflowexecution.DeleteWorkflowExecutionCommand
			mockBehavior func(*mocks.FakeWorkflowExecutionRepository)
			shouldErr    bool
		}{
			{
				name: "Should delete successfully",
				cmd: workflowexecution.DeleteWorkflowExecutionCommand{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeWorkflowExecutionRepository) {
					r.On("Delete", ctx, validID).Return(nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when validate command failed",
				cmd: workflowexecution.DeleteWorkflowExecutionCommand{
					ID: uuid.Nil,
				},
				mockBehavior: func(r *mocks.FakeWorkflowExecutionRepository) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository fails to delete",
				cmd: workflowexecution.DeleteWorkflowExecutionCommand{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeWorkflowExecutionRepository) {
					r.On("Delete", ctx, validID).Return(assert.AnError)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				repo := mocks.NewFakeWorkflowExecutionRepository(t)
				tc.mockBehavior(repo)

				s := workflowexecution.NewService(repo)
				err := s.Delete(ctx, tc.cmd)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("GetByID", func(t *testing.T) {
		testCases := []struct {
			name         string
			q            workflowexecution.GetWorkflowExecutionByIDQuery
			mockBehavior func(*mocks.FakeWorkflowExecutionRepository)
			shouldErr    bool
		}{
			{
				name: "Should get successfully",
				q: workflowexecution.GetWorkflowExecutionByIDQuery{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeWorkflowExecutionRepository) {
					r.On("Get", ctx, validID).Return(validWorkflowExecution, nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when validate query failed",
				q: workflowexecution.GetWorkflowExecutionByIDQuery{
					ID: uuid.Nil,
				},
				mockBehavior: func(r *mocks.FakeWorkflowExecutionRepository) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository fails to get",
				q: workflowexecution.GetWorkflowExecutionByIDQuery{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeWorkflowExecutionRepository) {
					r.On("Get", ctx, validID).Return(model.WorkflowExecution{}, assert.AnError)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				repo := mocks.NewFakeWorkflowExecutionRepository(t)
				tc.mockBehavior(repo)

				s := workflowexecution.NewService(repo)
				result, err := s.GetByID(ctx, tc.q)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.Equal(t, validWorkflowExecution, result)
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("GetStatusByID", func(t *testing.T) {
		testCases := []struct {
			name         string
			q            workflowexecution.GetWorkflowExecutionStatusByIDQuery
			mockBehavior func(*mocks.FakeWorkflowExecutionRepository)
			shouldErr    bool
		}{
			{
				name: "Should get status successfully",
				q: workflowexecution.GetWorkflowExecutionStatusByIDQuery{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeWorkflowExecutionRepository) {
					r.On("GetStatus", ctx, validID).Return(model.WorkflowExecutionStatusCompleted, nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when validate query failed",
				q: workflowexecution.GetWorkflowExecutionStatusByIDQuery{
					ID: uuid.Nil,
				},
				mockBehavior: func(r *mocks.FakeWorkflowExecutionRepository) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository fails to get status",
				q: workflowexecution.GetWorkflowExecutionStatusByIDQuery{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeWorkflowExecutionRepository) {
					r.On("GetStatus", ctx, validID).Return(model.WorkflowExecutionStatus(""), assert.AnError)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				repo := mocks.NewFakeWorkflowExecutionRepository(t)
				tc.mockBehavior(repo)

				s := workflowexecution.NewService(repo)
				result, err := s.GetStatusByID(ctx, tc.q)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.NotNil(t, result)
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("List", func(t *testing.T) {
		validPagingListWorkflowExecution := paging.List[model.WorkflowExecution]{
			Items:     []model.WorkflowExecution{validWorkflowExecution},
			TotalItem: 1,
		}

		testCases := []struct {
			name         string
			q            workflowexecution.ListWorkflowExecutionQuery
			mockBehavior func(*mocks.FakeWorkflowExecutionRepository)
			shouldErr    bool
		}{
			{
				name: "Should list successfully",
				q: workflowexecution.ListWorkflowExecutionQuery{
					WorkflowID:   validID,
					PagingParams: paging.NewParams(nil, nil),
					Sorts:        []xsort.Sort{},
				},
				mockBehavior: func(r *mocks.FakeWorkflowExecutionRepository) {
					r.On("List", ctx, validID, mock.Anything, mock.Anything).Return(&validPagingListWorkflowExecution, nil)
				},
				shouldErr: false,
			}, {
				name: "Should return error when validate query failed",
				q: workflowexecution.ListWorkflowExecutionQuery{
					WorkflowID:   uuid.Nil,
					PagingParams: paging.NewParams(nil, nil),
					Sorts:        []xsort.Sort{},
				},
				mockBehavior: func(r *mocks.FakeWorkflowExecutionRepository) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository fails to list",
				q: workflowexecution.ListWorkflowExecutionQuery{
					WorkflowID:   validID,
					PagingParams: paging.NewParams(nil, nil),
					Sorts:        []xsort.Sort{},
				},
				mockBehavior: func(r *mocks.FakeWorkflowExecutionRepository) {
					r.On("List", ctx, validID, mock.Anything, mock.Anything).Return(nil, assert.AnError)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				repo := mocks.NewFakeWorkflowExecutionRepository(t)
				tc.mockBehavior(repo)

				s := workflowexecution.NewService(repo)
				result, err := s.List(ctx, tc.q)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.Equal(t, validPagingListWorkflowExecution, *result)
					assert.NoError(t, err)
				}
			})
		}
	})
}

var (
	validID                = uuid.New()
	validWorkflowExecution = model.WorkflowExecution{
		ID:          validID,
		WorkflowID:  validID,
		Status:      model.WorkflowExecutionStatusCompleted,
		Env:         map[string]string{},
		Definition:  model.WorkflowDefinition{},
		CreatedAt:   time.Now(),
		StartedAt:   nil,
		CompletedAt: nil,
		Steps:       nil,
	}
)
