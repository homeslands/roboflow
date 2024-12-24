package step_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/internal/model/mocks"
	"github.com/tuanvumaihuynh/roboflow/internal/service/step"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

func TestStepService(t *testing.T) {
	ctx := context.Background()

	t.Run("GetByID", func(t *testing.T) {
		testCases := []struct {
			name         string
			query        step.GetStepByIDQuery
			mockBehavior func(*mocks.FakeStepRepository)
			shouldErr    bool
		}{
			{
				name: "Should get successfully",
				query: step.GetStepByIDQuery{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeStepRepository) {
					r.On("Get", ctx, mock.Anything).Return(validStep, nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when validate query failed",
				query: step.GetStepByIDQuery{
					ID: uuid.Nil,
				},
				mockBehavior: func(r *mocks.FakeStepRepository) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository return error",
				query: step.GetStepByIDQuery{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeStepRepository) {
					r.On("Get", ctx, mock.Anything).Return(model.Step{}, assert.AnError)
				},
				shouldErr: true,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				repo := mocks.NewFakeStepRepository(t)
				tc.mockBehavior(repo)

				svc := step.NewService(repo)
				step, err := svc.GetByID(ctx, tc.query)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.Equal(t, validStep, step)
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("List", func(t *testing.T) {
		testCases := []struct {
			name         string
			query        step.ListStepQuery
			mockBehavior func(*mocks.FakeStepRepository)
			shouldErr    bool
		}{
			{
				name: "Should list successfully",
				query: step.ListStepQuery{
					WorkflowExecutionID: validID,
					Sorts:               []xsort.Sort{},
				},
				mockBehavior: func(r *mocks.FakeStepRepository) {
					r.On("List", ctx, validID, mock.Anything, mock.Anything).Return([]model.Step{validStep}, nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when validate query failed",
				query: step.ListStepQuery{
					WorkflowExecutionID: uuid.Nil,
					Sorts:               []xsort.Sort{},
				},
				mockBehavior: func(r *mocks.FakeStepRepository) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository return error",
				query: step.ListStepQuery{
					WorkflowExecutionID: validID,
					Sorts:               []xsort.Sort{},
				},
				mockBehavior: func(r *mocks.FakeStepRepository) {
					r.On("List", ctx, validID, mock.Anything, mock.Anything).Return(nil, assert.AnError)
				},
				shouldErr: true,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				repo := mocks.NewFakeStepRepository(t)
				tc.mockBehavior(repo)

				svc := step.NewService(repo)
				result, err := svc.List(ctx, tc.query)

				if tc.shouldErr {
					assert.Nil(t, result)
					assert.Error(t, err)
				} else {
					assert.NotEmpty(t, result)
					assert.NoError(t, err)
				}
			})
		}
	})
}

var (
	validID   = uuid.New()
	validTime = time.Now()
	validStep = model.Step{
		ID:                  validID,
		WorkflowExecutionID: validID,
		Env:                 map[string]string{},
		Node:                model.WorkflowNode{},
		Status:              model.WorkflowExecutionStepStatusPending,
		StartedAt:           &validTime,
		CompletedAt:         &validTime,
	}
)
