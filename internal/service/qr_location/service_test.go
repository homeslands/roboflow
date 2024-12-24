package qrlocation_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/internal/model/mocks"
	qrlocation "github.com/tuanvumaihuynh/roboflow/internal/service/qr_location"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

func TestQRLocationService(t *testing.T) {
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		testCases := []struct {
			name         string
			cmd          qrlocation.CreateQrLocationCommand
			mockBehavior func(*mocks.FakeQRLocationRepository)
			shouldErr    bool
		}{
			{
				name: "Should create successfully",
				cmd: qrlocation.CreateQrLocationCommand{
					Name:     "Test",
					QRCode:   "Test",
					Metadata: map[string]interface{}{},
				},
				mockBehavior: func(r *mocks.FakeQRLocationRepository) {
					r.On("Create", ctx, mock.Anything).Return(nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when validate command failed",
				cmd: qrlocation.CreateQrLocationCommand{
					Name:     "",
					QRCode:   "Test",
					Metadata: map[string]interface{}{},
				},
				mockBehavior: func(r *mocks.FakeQRLocationRepository) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository fails to create",
				cmd: qrlocation.CreateQrLocationCommand{
					Name:     "Test",
					QRCode:   "Test",
					Metadata: map[string]interface{}{},
				},
				mockBehavior: func(r *mocks.FakeQRLocationRepository) {
					r.On("Create", ctx, mock.Anything).Return(assert.AnError)
				},
				shouldErr: true,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				qrLocationRepo := mocks.NewFakeQRLocationRepository(t)
				tc.mockBehavior(qrLocationRepo)

				qrLocationSvc := qrlocation.NewService(qrLocationRepo)
				result, err := qrLocationSvc.Create(ctx, tc.cmd)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.NotEmpty(t, result.ID)
					assert.Equal(t, tc.cmd.Name, result.Name)
					assert.Equal(t, tc.cmd.QRCode, result.QRCode)
					assert.Equal(t, tc.cmd.Metadata, result.Metadata)
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
			cmd          qrlocation.UpdateQRLocationCommand
			mockBehavior func(*mocks.FakeQRLocationRepository)
			shouldErr    bool
		}{
			{
				name: "Should update successfully",
				cmd: qrlocation.UpdateQRLocationCommand{
					ID:       validID,
					Name:     "Test",
					QRCode:   "Test",
					Metadata: map[string]interface{}{},
				},
				mockBehavior: func(r *mocks.FakeQRLocationRepository) {
					r.On("Update", ctx, mock.Anything).Return(validQrLoc, nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when validate command failed",
				cmd: qrlocation.UpdateQRLocationCommand{
					ID:       validID,
					Name:     "",
					QRCode:   "Test",
					Metadata: map[string]interface{}{},
				},
				mockBehavior: func(r *mocks.FakeQRLocationRepository) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository fails to update",
				cmd: qrlocation.UpdateQRLocationCommand{
					ID:       validID,
					Name:     "Test",
					QRCode:   "Test",
					Metadata: map[string]interface{}{},
				},
				mockBehavior: func(r *mocks.FakeQRLocationRepository) {
					r.On("Update", ctx, mock.Anything).Return(model.QRLocation{}, assert.AnError)
				},
				shouldErr: true,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				qrLocationRepo := mocks.NewFakeQRLocationRepository(t)
				tc.mockBehavior(qrLocationRepo)

				qrLocationSvc := qrlocation.NewService(qrLocationRepo)
				result, err := qrLocationSvc.Update(ctx, tc.cmd)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.Equal(t, tc.cmd.ID, result.ID)
					assert.Equal(t, tc.cmd.Name, result.Name)
					assert.Equal(t, tc.cmd.QRCode, result.QRCode)
					assert.Equal(t, tc.cmd.Metadata, result.Metadata)
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
			cmd          qrlocation.DeleteQrLocationCommand
			mockBehavior func(*mocks.FakeQRLocationRepository)
			shouldErr    bool
		}{
			{
				name: "Should delete successfully",
				cmd: qrlocation.DeleteQrLocationCommand{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeQRLocationRepository) {
					r.On("Delete", ctx, validID).Return(nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when validate command failed",
				cmd: qrlocation.DeleteQrLocationCommand{
					ID: uuid.Nil,
				},
				mockBehavior: func(r *mocks.FakeQRLocationRepository) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository fails to delete",
				cmd: qrlocation.DeleteQrLocationCommand{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeQRLocationRepository) {
					r.On("Delete", ctx, validID).Return(assert.AnError)
				},
				shouldErr: true,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				qrLocationRepo := mocks.NewFakeQRLocationRepository(t)
				tc.mockBehavior(qrLocationRepo)

				qrLocationSvc := qrlocation.NewService(qrLocationRepo)
				err := qrLocationSvc.Delete(ctx, tc.cmd)

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
			query        qrlocation.GetQrLocationByIDQuery
			mockBehavior func(*mocks.FakeQRLocationRepository)
			shouldErr    bool
		}{
			{
				name:  "Should get successfully",
				query: qrlocation.GetQrLocationByIDQuery{ID: validID},
				mockBehavior: func(r *mocks.FakeQRLocationRepository) {
					r.On("Get", ctx, validID).Return(validQrLoc, nil)
				},
				shouldErr: false,
			},
			{
				name:  "Should return error when validate query failed",
				query: qrlocation.GetQrLocationByIDQuery{ID: uuid.Nil},
				mockBehavior: func(r *mocks.FakeQRLocationRepository) {
				},
				shouldErr: true,
			},
			{
				name:  "Should return error when repository fails to get",
				query: qrlocation.GetQrLocationByIDQuery{ID: validID},
				mockBehavior: func(r *mocks.FakeQRLocationRepository) {
					r.On("Get", ctx, validID).Return(model.QRLocation{}, assert.AnError)
				},
				shouldErr: true,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				qrLocationRepo := mocks.NewFakeQRLocationRepository(t)
				tc.mockBehavior(qrLocationRepo)

				qrLocationSvc := qrlocation.NewService(qrLocationRepo)
				result, err := qrLocationSvc.GetByID(ctx, tc.query)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.Equal(t, validQrLoc, result)
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("List", func(t *testing.T) {
		validPagingListQrLocations := paging.List[model.QRLocation]{
			Items:     []model.QRLocation{validQrLoc},
			TotalItem: 1,
		}

		testCases := []struct {
			name         string
			query        qrlocation.ListQrLocationQuery
			mockBehavior func(*mocks.FakeQRLocationRepository)
			shouldErr    bool
		}{
			{
				name: "Should list successfully",
				query: qrlocation.ListQrLocationQuery{
					PagingParams: paging.NewParams(nil, nil),
					Sorts:        []xsort.Sort{},
				},
				mockBehavior: func(r *mocks.FakeQRLocationRepository) {
					r.On("List", ctx, mock.Anything, mock.Anything).Return(&validPagingListQrLocations, nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when validate command failed",
				query: qrlocation.ListQrLocationQuery{
					PagingParams: paging.NewParams(nil, nil),
					Sorts: []xsort.Sort{
						{
							Col:   "invalid_col",
							Order: xsort.OrderASC,
						},
					},
				},
				mockBehavior: func(r *mocks.FakeQRLocationRepository) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository fails to list",
				query: qrlocation.ListQrLocationQuery{
					PagingParams: paging.NewParams(nil, nil),
					Sorts:        []xsort.Sort{},
				},
				mockBehavior: func(r *mocks.FakeQRLocationRepository) {
					r.On("List", ctx, mock.Anything, mock.Anything).Return(nil, assert.AnError)
				},
				shouldErr: true,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				qrLocationRepo := mocks.NewFakeQRLocationRepository(t)
				tc.mockBehavior(qrLocationRepo)

				qrLocationSvc := qrlocation.NewService(qrLocationRepo)
				result, err := qrLocationSvc.List(ctx, tc.query)

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
	validID    = uuid.New()
	validQrLoc = model.QRLocation{
		ID:        validID,
		Name:      "Test",
		QRCode:    "Test",
		Metadata:  map[string]interface{}{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
)
