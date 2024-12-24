package raybot_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/internal/model/mocks"
	"github.com/tuanvumaihuynh/roboflow/internal/service/raybot"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

func TestRaybotService(t *testing.T) {
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		testCases := []struct {
			name         string
			cmd          raybot.CreateRaybotCommand
			mockBehavior func(*mocks.FakeRaybotRepository)
			shouldErr    bool
		}{
			{
				name: "Should create successfully",
				cmd: raybot.CreateRaybotCommand{
					Name: "Raybot",
				},
				mockBehavior: func(r *mocks.FakeRaybotRepository) {
					r.On("Create", ctx, mock.Anything).Return(nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when validate command failed",
				cmd: raybot.CreateRaybotCommand{
					Name: "",
				},
				mockBehavior: func(r *mocks.FakeRaybotRepository) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository fails to create",
				cmd: raybot.CreateRaybotCommand{
					Name: "Raybot",
				},
				mockBehavior: func(r *mocks.FakeRaybotRepository) {
					r.On("Create", ctx, mock.Anything).Return(assert.AnError)
				},
				shouldErr: true,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				raybotRepo := mocks.NewFakeRaybotRepository(t)
				tc.mockBehavior(raybotRepo)

				s := raybot.NewService(raybotRepo)
				result, err := s.Create(ctx, tc.cmd)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.NotEmpty(t, validRaybot.ID)
					assert.Equal(t, tc.cmd.Name, result.Name)
					assert.Equal(t, model.RaybotStatusOffline, result.Status)
					assert.NotEmpty(t, result.Token)
					assert.Empty(t, result.IpAddress)
					assert.Empty(t, result.LastConnectedAt)
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
			cmd          raybot.DeleteRaybotCommand
			mockBehavior func(*mocks.FakeRaybotRepository)
			shouldErr    bool
		}{
			{
				name: "Should delete successfully",
				cmd: raybot.DeleteRaybotCommand{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeRaybotRepository) {
					r.On("Delete", ctx, validID).Return(nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when validate command failed",
				cmd: raybot.DeleteRaybotCommand{
					ID: uuid.Nil,
				},
				mockBehavior: func(r *mocks.FakeRaybotRepository) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository fails to delete",
				cmd: raybot.DeleteRaybotCommand{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeRaybotRepository) {
					r.On("Delete", ctx, validID).Return(assert.AnError)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				raybotRepo := mocks.NewFakeRaybotRepository(t)
				tc.mockBehavior(raybotRepo)

				s := raybot.NewService(raybotRepo)
				err := s.Delete(ctx, tc.cmd)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("UpdateState", func(t *testing.T) {
		testCases := []struct {
			name         string
			cmd          raybot.UpdateStateCommand
			mockBehavior func(*mocks.FakeRaybotRepository)
			shouldErr    bool
		}{
			{
				name: "Should update state successfully",
				cmd: raybot.UpdateStateCommand{
					ID:    validID,
					State: model.RaybotStatusIdle,
				},
				mockBehavior: func(r *mocks.FakeRaybotRepository) {
					r.On("UpdateState", ctx, validID, model.RaybotStatusIdle).Return(nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when validate command failed",
				cmd: raybot.UpdateStateCommand{
					ID:    uuid.Nil,
					State: model.RaybotStatusIdle,
				},
				mockBehavior: func(r *mocks.FakeRaybotRepository) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository fails to update state",
				cmd: raybot.UpdateStateCommand{
					ID:    validID,
					State: model.RaybotStatusIdle,
				},
				mockBehavior: func(r *mocks.FakeRaybotRepository) {
					r.On("UpdateState", ctx, validID, model.RaybotStatusIdle).Return(assert.AnError)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				raybotRepo := mocks.NewFakeRaybotRepository(t)
				tc.mockBehavior(raybotRepo)

				s := raybot.NewService(raybotRepo)
				err := s.UpdateState(ctx, tc.cmd)

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
			query        raybot.GetRaybotByIDQuery
			mockBehavior func(*mocks.FakeRaybotRepository)
			shouldErr    bool
		}{
			{
				name: "Should get successfully",
				query: raybot.GetRaybotByIDQuery{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeRaybotRepository) {
					r.On("Get", ctx, validID).Return(validRaybot, nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when validate query failed",
				query: raybot.GetRaybotByIDQuery{
					ID: uuid.Nil,
				},
				mockBehavior: func(r *mocks.FakeRaybotRepository) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository fails to get",
				query: raybot.GetRaybotByIDQuery{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeRaybotRepository) {
					r.On("Get", ctx, validID).Return(model.Raybot{}, assert.AnError)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				raybotRepo := mocks.NewFakeRaybotRepository(t)
				tc.mockBehavior(raybotRepo)

				s := raybot.NewService(raybotRepo)
				result, err := s.GetByID(ctx, tc.query)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.Equal(t, validRaybot, result)
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("List", func(t *testing.T) {
		validPagingListRaybots := paging.List[model.Raybot]{
			Items:     []model.Raybot{validRaybot, validRaybot, validRaybot},
			TotalItem: 3,
		}

		testCases := []struct {
			name         string
			query        raybot.ListRaybotQuery
			mockBehavior func(*mocks.FakeRaybotRepository)
			shouldErr    bool
		}{
			{
				name: "Should list successfully",
				query: raybot.ListRaybotQuery{
					PagingParams: paging.NewParams(nil, nil),
					Sorts:        []xsort.Sort{},
				},
				mockBehavior: func(r *mocks.FakeRaybotRepository) {
					r.On("List", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&validPagingListRaybots, nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when validate command failed",
				query: raybot.ListRaybotQuery{
					PagingParams: paging.NewParams(nil, nil),
					Sorts: []xsort.Sort{
						{
							Col:   "invalid_col",
							Order: xsort.OrderASC,
						},
					},
				},
				mockBehavior: func(r *mocks.FakeRaybotRepository) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository fails to list",
				query: raybot.ListRaybotQuery{
					PagingParams: paging.NewParams(nil, nil),
					Sorts:        []xsort.Sort{},
				},
				mockBehavior: func(r *mocks.FakeRaybotRepository) {
					r.On("List", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				raybotRepo := mocks.NewFakeRaybotRepository(t)
				tc.mockBehavior(raybotRepo)

				s := raybot.NewService(raybotRepo)
				result, err := s.List(ctx, tc.query)

				if tc.shouldErr {
					assert.Nil(t, result)
					assert.Error(t, err)
				} else {
					assert.NotEmpty(t, result.Items)
					assert.Equal(t, validPagingListRaybots.TotalItem, result.TotalItem)
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("GetState", func(t *testing.T) {
		testCases := []struct {
			name         string
			query        raybot.GetStatusQuery
			mockBehavior func(*mocks.FakeRaybotRepository)
			shouldErr    bool
		}{
			{
				name: "Should get state successfully",
				query: raybot.GetStatusQuery{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeRaybotRepository) {
					r.On("GetState", ctx, validID).Return(model.RaybotStatusIdle, nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when validate command failed",
				query: raybot.GetStatusQuery{
					ID: uuid.Nil,
				},
				mockBehavior: func(r *mocks.FakeRaybotRepository) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository fails to get state",
				query: raybot.GetStatusQuery{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeRaybotRepository) {
					r.On("GetState", ctx, validID).Return(model.RaybotStatus(""), assert.AnError)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				raybotRepo := mocks.NewFakeRaybotRepository(t)
				tc.mockBehavior(raybotRepo)

				s := raybot.NewService(raybotRepo)
				result, err := s.GetState(ctx, tc.query)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.Equal(t, model.RaybotStatusIdle, result)
					assert.NoError(t, err)
				}
			})
		}
	})
}

var (
	validID     = uuid.New()
	validRaybot = model.Raybot{
		ID:              validID,
		Name:            "Raybot",
		Token:           uuid.NewString(),
		Status:          model.RaybotStatusOffline,
		IpAddress:       nil,
		LastConnectedAt: nil,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
)
