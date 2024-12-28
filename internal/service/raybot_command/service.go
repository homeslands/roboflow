package raybotcommand

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
	"github.com/tuanvumaihuynh/roboflow/internal/service/raybot_command/event"
	"github.com/tuanvumaihuynh/roboflow/pkg/paging"
	"github.com/tuanvumaihuynh/roboflow/pkg/pubsub"
	"github.com/tuanvumaihuynh/roboflow/pkg/validator"
	"github.com/tuanvumaihuynh/roboflow/pkg/xerrors"
)

var _ Service = (*service)(nil)

type Service interface {
	Create(ctx context.Context, cmd CreateRaybotCommandCommand) (model.RaybotCommand, error)
	Delete(ctx context.Context, cmd DeleteRaybotCommandCommand) error
	SetStatusInProgess(ctx context.Context, cmd SetStatusInProgessCommand) error
	SetStatusSuccess(ctx context.Context, cmd SetStatusSuccessCommand) error
	SetStatusFailed(ctx context.Context, cmd SetStatusFailedCommand) error

	GetByID(ctx context.Context, q GetRaybotCommandByIDQuery) (model.RaybotCommand, error)
	List(ctx context.Context, q ListRaybotCommandQuery) (*paging.List[model.RaybotCommand], error)
}

type service struct {
	raybotCommandRepo model.RaybotCommandRepository
	raybotRepo        model.RaybotRepository
	qrLoccationRepo   model.QRLocationRepository
	eventPublisher    pubsub.Publisher
	log               *slog.Logger
}

func NewService(
	raybotCommandRepo model.RaybotCommandRepository,
	raybotRepo model.RaybotRepository,
	qrLoccationRepo model.QRLocationRepository,
	eventPublisher pubsub.Publisher,
	log *slog.Logger,
) *service {
	return &service{
		raybotCommandRepo: raybotCommandRepo,
		raybotRepo:        raybotRepo,
		qrLoccationRepo:   qrLoccationRepo,
		eventPublisher:    eventPublisher,
		log:               log,
	}
}

func (s service) Create(ctx context.Context, cmd CreateRaybotCommandCommand) (model.RaybotCommand, error) {
	if err := cmd.Validate(); err != nil {
		return model.RaybotCommand{}, err
	}

	// Validate the state of the raybot before creating the command
	if err := s.validateCommandState(ctx, cmd); err != nil {
		return model.RaybotCommand{}, err
	}

	// Build command input
	err := validateCommandInput(cmd.Type, cmd.Input)
	if err != nil {
		return model.RaybotCommand{}, err
	}

	// Create raybot command
	modelRaybotCommand := model.RaybotCommand{
		ID:        uuid.New(),
		RaybotID:  cmd.RaybotID,
		Type:      cmd.Type,
		Status:    model.RaybotCommandStatusPending,
		Input:     cmd.Input,
		CreatedAt: time.Now(),
	}
	if err := s.raybotCommandRepo.Create(ctx, modelRaybotCommand); err != nil {
		return model.RaybotCommand{}, err
	}

	// Publish event
	if err := s.publishRaybotCommandCreatedEvent(modelRaybotCommand); err != nil {
		return model.RaybotCommand{}, err
	}

	return modelRaybotCommand, nil
}

func (s service) Delete(ctx context.Context, cmd DeleteRaybotCommandCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	return s.raybotCommandRepo.Delete(ctx, cmd.ID)
}

func (s service) SetStatusInProgess(ctx context.Context, cmd SetStatusInProgessCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	// When raybot receives the command in progress, raybot stautus will be BUSY
	return s.raybotCommandRepo.Update(
		ctx,
		cmd.ID,
		model.RaybotStatusBusy,
		func(raybotCmd *model.RaybotCommand) error {
			if raybotCmd.Status != model.RaybotCommandStatusPending {
				return xerrors.ThrowPreconditionFailed(nil, "command status is not PENDING")
			}
			raybotCmd.Status = model.RaybotCommandStatusInProgress

			return nil
		},
	)
}

func (s service) SetStatusSuccess(ctx context.Context, cmd SetStatusSuccessCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	// When raybot completes the command, raybot status will be IDLE
	return s.raybotCommandRepo.Update(
		ctx,
		cmd.ID,
		model.RaybotStatusIdle,
		func(raybotCmd *model.RaybotCommand) error {
			if raybotCmd.Status != model.RaybotCommandStatusInProgress {
				return xerrors.ThrowPreconditionFailed(nil, "command status is not IN_PROGRESS")
			}
			raybotCmd.Status = model.RaybotCommandStatusSuccess
			now := time.Now()
			raybotCmd.CompletedAt = &now
			raybotCmd.Output = cmd.Output

			return nil
		},
	)
}

func (s service) SetStatusFailed(ctx context.Context, cmd SetStatusFailedCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	// When raybot fails to complete the command, raybot status will be IDLE
	return s.raybotCommandRepo.Update(
		ctx,
		cmd.ID,
		model.RaybotStatusIdle,
		func(raybotCmd *model.RaybotCommand) error {
			if raybotCmd.Status != model.RaybotCommandStatusInProgress &&
				raybotCmd.Status != model.RaybotCommandStatusPending {
				return xerrors.ThrowPreconditionFailed(nil, "command status is not IN_PROGRESS or PENDING")
			}
			raybotCmd.Status = model.RaybotCommandStatusFailed
			now := time.Now()
			raybotCmd.CompletedAt = &now
			raybotCmd.Output = cmd.Output

			return nil
		},
	)
}

func (s service) GetByID(ctx context.Context, q GetRaybotCommandByIDQuery) (model.RaybotCommand, error) {
	if err := q.Validate(); err != nil {
		return model.RaybotCommand{}, err
	}

	return s.raybotCommandRepo.Get(ctx, q.ID)
}

func (s service) List(ctx context.Context, q ListRaybotCommandQuery) (*paging.List[model.RaybotCommand], error) {
	if err := q.Validate(); err != nil {
		return &paging.List[model.RaybotCommand]{}, err
	}

	return s.raybotCommandRepo.List(ctx, q.RaybotID, q.PagingParams, q.Sorts)
}

// Validate the state of the raybot
func (s service) validateCommandState(ctx context.Context, cmd CreateRaybotCommandCommand) error {
	state, err := s.raybotRepo.GetState(ctx, cmd.RaybotID)
	if err != nil {
		return err
	}
	// Raybot must be ONLINE to receive command
	if state == model.RaybotStatusOffline {
		return xerrors.ThrowPreconditionFailed(nil, "raybot is not ONLINE")
	}
	// Only STOP command is allowed when raybot is BUSY
	if state == model.RaybotStatusBusy && cmd.Type != model.RaybotCommandTypeStop {
		return xerrors.ThrowPreconditionFailed(nil, "raybot is BUSY, only STOP command is allowed")
	}
	return nil
}

// Build the appropriate input for the command type
func validateCommandInput(cmdType model.RaybotCommandType, input []byte) error {
	switch cmdType {
	case model.RaybotCommandTypeMoveToLocation:
		var i model.MoveToLocationInput
		if err := json.Unmarshal(input, &i); err != nil {
			return xerrors.ThrowInvalidArgument(err, "invalid input")
		}
		return validator.Validate(i)

	case model.RaybotCommandTypeCheckQrCode:
		var i model.CheckQRCodeInput
		if err := json.Unmarshal(input, &i); err != nil {
			return xerrors.ThrowInvalidArgument(err, "invalid input")
		}
		return validator.Validate(i)

	case model.RaybotCommandTypeLiftBox, model.RaybotCommandTypeDropBox:
		var i model.LiftDropBoxInput
		if err := json.Unmarshal(input, &i); err != nil {
			return xerrors.ThrowInvalidArgument(err, "invalid input")
		}
		return validator.Validate(i)

	case model.RaybotCommandTypeStop,
		model.RaybotCommandTypeMoveForward,
		model.RaybotCommandTypeMoveBackward,
		model.RaybotCommandTypeOpenBox,
		model.RaybotCommandTypeCloseBox,
		model.RaybotCommandTypeWaitGetItem:
		// No input required
		if len(input) != 0 {
			return xerrors.ThrowInvalidArgument(nil, "input is not required for this command type")
		}

		return nil

	default:
		return xerrors.ThrowInvalidArgument(nil, "unknown command type")
	}
}

func (s service) publishRaybotCommandCreatedEvent(modelRaybotCommand model.RaybotCommand) error {
	ev := event.RaybotCommandCreated(modelRaybotCommand)
	msgJSON, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	if err := s.eventPublisher.Publish(event.TopicRaybotCommandCreated, msgJSON); err != nil {
		return err
	}

	return nil
}
