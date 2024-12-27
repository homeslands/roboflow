package raybotcommand_test

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
	raybotcommand "github.com/tuanvumaihuynh/roboflow/internal/service/raybot_command"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	pubsubmocks "github.com/tuanvumaihuynh/roboflow/pkg/pubsub/mocks"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
	"github.com/tuanvumaihuynh/roboflow/pkg/xsort"
)

func TestRaybotCommandService(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		type testCase struct {
			name              string
			cmd               raybotcommand.CreateRaybotCommandCommand
			raybotCommandRepo *mocks.FakeRaybotCommandRepository
			raybotRepo        *mocks.FakeRaybotRepository
			qrLocationRepo    *mocks.FakeQRLocationRepository
			eventPublisher    *pubsubmocks.FakePublisher
			mockBehavior      func(*testCase)
			shouldErr         bool
		}

		testCases := []testCase{
			{
				name: "Should create successfully",
				cmd: raybotcommand.CreateRaybotCommandCommand{
					RaybotID: validID,
					Type:     model.RaybotCommandTypeMoveForward,
					Input:    struct{}{},
				},
				raybotCommandRepo: mocks.NewFakeRaybotCommandRepository(t),
				raybotRepo:        mocks.NewFakeRaybotRepository(t),
				qrLocationRepo:    mocks.NewFakeQRLocationRepository(t),
				eventPublisher:    pubsubmocks.NewFakePublisher(t),
				mockBehavior: func(tc *testCase) {
					tc.raybotRepo.EXPECT().GetState(ctx, validID).Return(model.RaybotStatusIdle, nil)
					tc.raybotCommandRepo.EXPECT().Create(ctx, mock.Anything).Return(nil)
					tc.eventPublisher.EXPECT().Publish("raybot.command.created", mock.Anything).Return(nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when validate command failed",
				cmd: raybotcommand.CreateRaybotCommandCommand{
					RaybotID: uuid.Nil,
					Type:     model.RaybotCommandTypeMoveForward,
					Input:    map[string]interface{}{},
				},
				raybotCommandRepo: mocks.NewFakeRaybotCommandRepository(t),
				raybotRepo:        mocks.NewFakeRaybotRepository(t),
				qrLocationRepo:    mocks.NewFakeQRLocationRepository(t),
				eventPublisher:    pubsubmocks.NewFakePublisher(t),
				mockBehavior: func(tc *testCase) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when get state fails",
				cmd: raybotcommand.CreateRaybotCommandCommand{
					RaybotID: validID,
					Type:     model.RaybotCommandTypeMoveForward,
					Input:    map[string]interface{}{},
				},
				raybotCommandRepo: mocks.NewFakeRaybotCommandRepository(t),
				raybotRepo:        mocks.NewFakeRaybotRepository(t),
				qrLocationRepo:    mocks.NewFakeQRLocationRepository(t),
				eventPublisher:    pubsubmocks.NewFakePublisher(t),
				mockBehavior: func(tc *testCase) {
					tc.raybotRepo.EXPECT().GetState(ctx, validID).Return(model.RaybotStatus(""), assert.AnError)
				},
				shouldErr: true,
			},
			{
				name: "Should return error when raybot is OFFLINE",
				cmd: raybotcommand.CreateRaybotCommandCommand{
					RaybotID: validID,
					Type:     model.RaybotCommandTypeMoveForward,
					Input:    map[string]interface{}{},
				},
				raybotCommandRepo: mocks.NewFakeRaybotCommandRepository(t),
				raybotRepo:        mocks.NewFakeRaybotRepository(t),
				qrLocationRepo:    mocks.NewFakeQRLocationRepository(t),
				eventPublisher:    pubsubmocks.NewFakePublisher(t),
				mockBehavior: func(tc *testCase) {
					tc.raybotRepo.EXPECT().GetState(ctx, validID).Return(model.RaybotStatusOffline, nil)
				},
				shouldErr: true,
			},
			{
				name: "Should return error when raybot is BUSY and command is not STOP",
				cmd: raybotcommand.CreateRaybotCommandCommand{
					RaybotID: validID,
					Type:     model.RaybotCommandTypeMoveForward,
					Input:    map[string]interface{}{},
				},
				raybotCommandRepo: mocks.NewFakeRaybotCommandRepository(t),
				raybotRepo:        mocks.NewFakeRaybotRepository(t),
				qrLocationRepo:    mocks.NewFakeQRLocationRepository(t),
				eventPublisher:    pubsubmocks.NewFakePublisher(t),
				mockBehavior: func(tc *testCase) {
					tc.raybotRepo.EXPECT().GetState(ctx, validID).Return(model.RaybotStatusBusy, nil)
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository fails to create",
				cmd: raybotcommand.CreateRaybotCommandCommand{
					RaybotID: validID,
					Type:     model.RaybotCommandTypeMoveForward,
					Input:    map[string]interface{}{},
				},
				raybotCommandRepo: mocks.NewFakeRaybotCommandRepository(t),
				raybotRepo:        mocks.NewFakeRaybotRepository(t),
				qrLocationRepo:    mocks.NewFakeQRLocationRepository(t),
				eventPublisher:    pubsubmocks.NewFakePublisher(t),
				mockBehavior: func(tc *testCase) {
					tc.raybotRepo.EXPECT().GetState(ctx, validID).Return(model.RaybotStatusIdle, nil)
					tc.raybotCommandRepo.EXPECT().Create(ctx, mock.Anything).Return(assert.AnError)
				},
				shouldErr: true,
			},
			{
				name: "Should return error when event publisher fails to publish",
				cmd: raybotcommand.CreateRaybotCommandCommand{
					RaybotID: validID,
					Type:     model.RaybotCommandTypeMoveForward,
					Input:    map[string]interface{}{},
				},
				raybotCommandRepo: mocks.NewFakeRaybotCommandRepository(t),
				raybotRepo:        mocks.NewFakeRaybotRepository(t),
				qrLocationRepo:    mocks.NewFakeQRLocationRepository(t),
				eventPublisher:    pubsubmocks.NewFakePublisher(t),
				mockBehavior: func(tc *testCase) {
					tc.raybotRepo.EXPECT().GetState(ctx, validID).Return(model.RaybotStatusIdle, nil)
					tc.raybotCommandRepo.EXPECT().Create(ctx, mock.Anything).Return(nil)
					tc.eventPublisher.EXPECT().Publish("raybot.command.created", mock.Anything).Return(assert.AnError)
				},
				shouldErr: true,
			},
			{
				name: "Should check QR code if command type is MOVE_TO_LOCATION",
				cmd: raybotcommand.CreateRaybotCommandCommand{
					RaybotID: validID,
					Type:     model.RaybotCommandTypeMoveToLocation,
					Input: map[string]interface{}{
						"location":  "qr_code",
						"direction": "FORWARD",
					},
				},
				raybotCommandRepo: mocks.NewFakeRaybotCommandRepository(t),
				raybotRepo:        mocks.NewFakeRaybotRepository(t),
				qrLocationRepo:    mocks.NewFakeQRLocationRepository(t),
				eventPublisher:    pubsubmocks.NewFakePublisher(t),
				mockBehavior: func(tc *testCase) {
					tc.raybotRepo.EXPECT().GetState(ctx, validID).Return(model.RaybotStatusIdle, nil)
					tc.raybotCommandRepo.EXPECT().Create(ctx, mock.Anything).Return(nil)
					tc.qrLocationRepo.EXPECT().ExistByQRCode(ctx, "qr_code").Return(true, nil)
					tc.eventPublisher.EXPECT().Publish("raybot.command.created", mock.Anything).Return(nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when check QR code fails",
				cmd: raybotcommand.CreateRaybotCommandCommand{
					RaybotID: validID,
					Type:     model.RaybotCommandTypeMoveToLocation,
					Input: map[string]interface{}{
						"location":  "qr_code",
						"direction": "FORWARD",
					},
				},
				raybotCommandRepo: mocks.NewFakeRaybotCommandRepository(t),
				raybotRepo:        mocks.NewFakeRaybotRepository(t),
				qrLocationRepo:    mocks.NewFakeQRLocationRepository(t),
				eventPublisher:    pubsubmocks.NewFakePublisher(t),
				mockBehavior: func(tc *testCase) {
					tc.raybotRepo.EXPECT().GetState(ctx, validID).Return(model.RaybotStatusIdle, nil)
					tc.qrLocationRepo.EXPECT().ExistByQRCode(ctx, "qr_code").Return(false, assert.AnError)
				},
				shouldErr: true,
			},
			{
				name: "Should return error when QR code does not exist",
				cmd: raybotcommand.CreateRaybotCommandCommand{
					RaybotID: validID,
					Type:     model.RaybotCommandTypeMoveToLocation,
					Input: map[string]interface{}{
						"location":  "qr_code",
						"direction": "FORWARD",
					},
				},
				raybotCommandRepo: mocks.NewFakeRaybotCommandRepository(t),
				raybotRepo:        mocks.NewFakeRaybotRepository(t),
				qrLocationRepo:    mocks.NewFakeQRLocationRepository(t),
				eventPublisher:    pubsubmocks.NewFakePublisher(t),
				mockBehavior: func(tc *testCase) {
					tc.raybotRepo.EXPECT().GetState(ctx, validID).Return(model.RaybotStatusIdle, nil)
					tc.qrLocationRepo.EXPECT().ExistByQRCode(ctx, "qr_code").Return(false, nil)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tc.mockBehavior(&tc)

				s := raybotcommand.NewService(tc.raybotCommandRepo, tc.raybotRepo, tc.qrLocationRepo, tc.eventPublisher, log)
				result, err := s.Create(ctx, tc.cmd)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.Equal(t, tc.cmd.RaybotID, result.RaybotID)
					assert.Equal(t, tc.cmd.Type, result.Type)
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("Delete", func(t *testing.T) {
		testCases := []struct {
			name         string
			cmd          raybotcommand.DeleteRaybotCommandCommand
			mockBehavior func(*mocks.FakeRaybotCommandRepository)
			shouldErr    bool
		}{
			{
				name: "Should delete successfully",
				cmd: raybotcommand.DeleteRaybotCommandCommand{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeRaybotCommandRepository) {
					r.On("Delete", ctx, validID).Return(nil)
				},
				shouldErr: false,
			}, {
				name: "Should return error when validate command failed",
				cmd: raybotcommand.DeleteRaybotCommandCommand{
					ID: uuid.Nil,
				},
				mockBehavior: func(r *mocks.FakeRaybotCommandRepository) {
				},
				shouldErr: true,
			}, {
				name: "Should return error when repository fails to delete",
				cmd: raybotcommand.DeleteRaybotCommandCommand{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeRaybotCommandRepository) {
					r.On("Delete", ctx, validID).Return(assert.AnError)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				repo := &mocks.FakeRaybotCommandRepository{}
				tc.mockBehavior(repo)

				s := raybotcommand.NewService(repo, nil, nil, nil, log)
				err := s.Delete(ctx, tc.cmd)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("SetStatusInProgress", func(t *testing.T) {
		testCases := []struct {
			name         string
			cmd          raybotcommand.SetStatusInProgessCommand
			mockBehavior func(*mocks.FakeRaybotCommandRepository)
			shoudErr     bool
		}{
			{
				name: "Should set status in progress successfully",
				cmd: raybotcommand.SetStatusInProgessCommand{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeRaybotCommandRepository) {
					r.EXPECT().Update(
						ctx,
						validID,
						model.RaybotStatusBusy,
						mock.MatchedBy(
							func(fn func(*model.RaybotCommand) error) bool {
								cmd := model.RaybotCommand{
									Status: model.RaybotCommandStatusPending,
								}
								err := fn(&cmd)
								return err == nil && cmd.Status == model.RaybotCommandStatusInProgress
							}),
					).Return(nil)
				},
				shoudErr: false,
			},
			{
				name: "Should return error when validate command failed",
				cmd: raybotcommand.SetStatusInProgessCommand{
					ID: uuid.Nil,
				},
				mockBehavior: func(r *mocks.FakeRaybotCommandRepository) {
				},
				shoudErr: true,
			},
			{
				name: "Should return error when current command status is not PENDING",
				cmd: raybotcommand.SetStatusInProgessCommand{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeRaybotCommandRepository) {
					r.EXPECT().Update(
						ctx,
						validID,
						model.RaybotStatusBusy,
						mock.MatchedBy(
							func(fn func(*model.RaybotCommand) error) bool {
								cmd := model.RaybotCommand{
									// Mock current command status is not PENDING
									Status: model.RaybotCommandStatusInProgress,
								}
								err := fn(&cmd)
								return err != nil && xerrors.IsPreconditionFailed(err)
							}),
					).Return(assert.AnError)
				},
				shoudErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				repo := mocks.NewFakeRaybotCommandRepository(t)
				tc.mockBehavior(repo)

				s := raybotcommand.NewService(repo, nil, nil, nil, log)
				err := s.SetStatusInProgess(ctx, tc.cmd)

				if tc.shoudErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("SetStatusSuccess", func(t *testing.T) {
		testCases := []struct {
			name         string
			cmd          raybotcommand.SetStatusSuccessCommand
			mockBehavior func(*mocks.FakeRaybotCommandRepository)
			shoudErr     bool
		}{
			{
				name: "Should set status success successfully",
				cmd: raybotcommand.SetStatusSuccessCommand{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeRaybotCommandRepository) {
					r.EXPECT().Update(
						ctx,
						validID,
						model.RaybotStatusIdle,
						mock.MatchedBy(
							func(fn func(*model.RaybotCommand) error) bool {
								cmd := model.RaybotCommand{
									Status: model.RaybotCommandStatusInProgress,
								}
								err := fn(&cmd)
								return err == nil && cmd.Status == model.RaybotCommandStatusSuccess
							}),
					).Return(nil)
				},
				shoudErr: false,
			},
			{
				name: "Should return error when validate command failed",
				cmd: raybotcommand.SetStatusSuccessCommand{
					ID: uuid.Nil,
				},
				mockBehavior: func(r *mocks.FakeRaybotCommandRepository) {
				},
				shoudErr: true,
			},
			{
				name: "Should return error when current command status is not IN_PROGRESS",
				cmd: raybotcommand.SetStatusSuccessCommand{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeRaybotCommandRepository) {
					r.EXPECT().Update(
						ctx,
						validID,
						model.RaybotStatusIdle,
						mock.MatchedBy(
							func(fn func(*model.RaybotCommand) error) bool {
								cmd := model.RaybotCommand{
									// Mock current command status is not IN_PROGRESS
									Status: model.RaybotCommandStatusPending,
								}
								err := fn(&cmd)
								return err != nil && xerrors.IsPreconditionFailed(err)
							}),
					).Return(assert.AnError)
				},
				shoudErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				repo := mocks.NewFakeRaybotCommandRepository(t)
				tc.mockBehavior(repo)

				s := raybotcommand.NewService(repo, nil, nil, nil, log)
				err := s.SetStatusSuccess(ctx, tc.cmd)

				if tc.shoudErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("SetStatusFailed", func(t *testing.T) {
		testCases := []struct {
			name         string
			cmd          raybotcommand.SetStatusFailedCommand
			mockBehavior func(*mocks.FakeRaybotCommandRepository)
			shoudErr     bool
		}{
			{
				name: "Should set status failed successfully",
				cmd: raybotcommand.SetStatusFailedCommand{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeRaybotCommandRepository) {
					r.EXPECT().Update(
						ctx,
						validID,
						model.RaybotStatusIdle,
						mock.MatchedBy(
							func(fn func(*model.RaybotCommand) error) bool {
								cmd := model.RaybotCommand{
									Status: model.RaybotCommandStatusInProgress,
								}
								err := fn(&cmd)
								return err == nil && cmd.Status == model.RaybotCommandStatusFailed
							}),
					).Return(nil)
				},
				shoudErr: false,
			},
			{
				name: "Should return error when validate command failed",
				cmd: raybotcommand.SetStatusFailedCommand{
					ID: uuid.Nil,
				},
				mockBehavior: func(r *mocks.FakeRaybotCommandRepository) {
				},
				shoudErr: true,
			},
			{
				name: "Should return error when current command status is not IN_PROGRESS",
				cmd: raybotcommand.SetStatusFailedCommand{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeRaybotCommandRepository) {
					r.EXPECT().Update(
						ctx,
						validID,
						model.RaybotStatusIdle,
						mock.MatchedBy(
							func(fn func(*model.RaybotCommand) error) bool {
								cmd := model.RaybotCommand{
									// Mock current command status is not IN_PROGRESS
									Status: model.RaybotCommandStatusPending,
								}
								err := fn(&cmd)
								return err != nil && xerrors.IsPreconditionFailed(err)
							}),
					).Return(assert.AnError)
				},
				shoudErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				repo := mocks.NewFakeRaybotCommandRepository(t)
				tc.mockBehavior(repo)

				s := raybotcommand.NewService(repo, nil, nil, nil, log)
				err := s.SetStatusFailed(ctx, tc.cmd)

				if tc.shoudErr {
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
			q            raybotcommand.GetRaybotCommandByIDQuery
			mockBehavior func(*mocks.FakeRaybotCommandRepository)
			shouldErr    bool
		}{
			{
				name: "Should get successfully",
				q: raybotcommand.GetRaybotCommandByIDQuery{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeRaybotCommandRepository) {
					r.On("Get", ctx, validID).Return(validRaybotCommand, nil)
				},
				shouldErr: false,
			}, {
				name: "Should return error when validate query failed",
				q: raybotcommand.GetRaybotCommandByIDQuery{
					ID: uuid.Nil,
				},
				mockBehavior: func(r *mocks.FakeRaybotCommandRepository) {
				},
				shouldErr: true,
			}, {
				name: "Should return error when repository fails to get",
				q: raybotcommand.GetRaybotCommandByIDQuery{
					ID: validID,
				},
				mockBehavior: func(r *mocks.FakeRaybotCommandRepository) {
					r.On("Get", ctx, validID).Return(validRaybotCommand, assert.AnError)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				raybotCommandRepo := mocks.NewFakeRaybotCommandRepository(t)
				tc.mockBehavior(raybotCommandRepo)

				s := raybotcommand.NewService(raybotCommandRepo, nil, nil, nil, log)
				_, err := s.GetByID(ctx, tc.q)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("List", func(t *testing.T) {
		validPagingListRaybotCommands := paging.List[model.RaybotCommand]{
			Items:     []model.RaybotCommand{validRaybotCommand},
			TotalItem: 1,
		}

		testCases := []struct {
			name         string
			q            raybotcommand.ListRaybotCommandQuery
			mockBehavior func(*mocks.FakeRaybotCommandRepository)
			shouldErr    bool
		}{
			{
				name: "Should list successfully",
				q: raybotcommand.ListRaybotCommandQuery{
					RaybotID:     validID,
					PagingParams: paging.NewParams(nil, nil),
					Sorts:        []xsort.Sort{},
				},
				mockBehavior: func(r *mocks.FakeRaybotCommandRepository) {
					r.On("List", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&validPagingListRaybotCommands, nil)
				},
				shouldErr: false,
			},
			{
				name: "Should return error when validate query failed",
				q: raybotcommand.ListRaybotCommandQuery{
					RaybotID:     uuid.Nil,
					PagingParams: paging.NewParams(nil, nil),
					Sorts:        []xsort.Sort{},
				},
				mockBehavior: func(r *mocks.FakeRaybotCommandRepository) {
				},
				shouldErr: true,
			},
			{
				name: "Should return error when repository fails to list",
				q: raybotcommand.ListRaybotCommandQuery{
					RaybotID:     validID,
					PagingParams: paging.NewParams(nil, nil),
					Sorts:        []xsort.Sort{},
				},
				mockBehavior: func(r *mocks.FakeRaybotCommandRepository) {
					r.On("List", ctx, validID, mock.Anything, mock.Anything).Return(&validPagingListRaybotCommands, assert.AnError)
				},
				shouldErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				raybotCommandRepo := mocks.NewFakeRaybotCommandRepository(t)
				tc.mockBehavior(raybotCommandRepo)

				s := raybotcommand.NewService(raybotCommandRepo, nil, nil, nil, log)
				_, err := s.List(ctx, tc.q)

				if tc.shouldErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

var (
	validID            = uuid.New()
	validRaybotCommand = model.RaybotCommand{
		RaybotID:    validID,
		ID:          validID,
		Type:        model.RaybotCommandTypeCheckQrCode,
		Status:      model.RaybotCommandStatusPending,
		Input:       map[string]interface{}{},
		Output:      map[string]interface{}{},
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}
)
